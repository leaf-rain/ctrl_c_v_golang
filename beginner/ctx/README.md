# 数组

***
ctrl c/v永不过时 ！！！好用请给star。

### 1. 应用场景
在并发程序中，由于超时、取消操作或者一些异常情况，往往需要进行抢占操作或者中断后续操作。
```go
func main() {
    msg := make(chan int, 5)
    ch := make(chan bool)

    defer close(msg)
    // consumer
    go func() {
		defer fmt.Println("consumer协程退出")
        ticker := time.NewTicker(1 * time.Second)
        for _ = range ticker.C {
            select {
            case <-ch:
                fmt.Println("收到退出信号")
                return
            default:
                fmt.Printf("收到消息: %d\n", <-msg)
            }
        }
    }()

    // producer
    for i := 0; i < 10; i++ {
        msg <- i
    }
    time.Sleep(5 * time.Second)
    close(ch)
    time.Sleep(1 * time.Second)
    fmt.Println("主协程退出")
}
```

### 2. [结构](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/context/context.go#L62)
我一直觉得看源代码上面的注释是最准确的(如果英语好的话...)
```go
// A Context carries a deadline, a cancellation signal, and other values across
// API boundaries.
//
// Context's methods may be called by multiple goroutines simultaneously.
type Context interface {
	// Deadline returns the time when work done on behalf of this context
	// should be canceled. Deadline returns ok==false when no deadline is
	// set. Successive calls to Deadline return the same results.
	Deadline() (deadline time.Time, ok bool)

	// Done returns a channel that's closed when work done on behalf of this
	// context should be canceled. Done may return nil if this context can
	// never be canceled. Successive calls to Done return the same value.
	// The close of the Done channel may happen asynchronously,
	// after the cancel function returns.
	//
	// WithCancel arranges for Done to be closed when cancel is called;
	// WithDeadline arranges for Done to be closed when the deadline
	// expires; WithTimeout arranges for Done to be closed when the timeout
	// elapses.
	//
	// Done is provided for use in select statements:
	//
	//  // Stream generates values with DoSomething and sends them to out
	//  // until DoSomething returns an error or ctx.Done is closed.
	//  func Stream(ctx context.Context, out chan<- Value) error {
	//  	for {
	//  		v, err := DoSomething(ctx)
	//  		if err != nil {
	//  			return err
	//  		}
	//  		select {
	//  		case <-ctx.Done():
	//  			return ctx.Err()
	//  		case out <- v:
	//  		}
	//  	}
	//  }
	//
	// See https://blog.golang.org/pipelines for more examples of how to use
	// a Done channel for cancellation.
	Done() <-chan struct{}

	// If Done is not yet closed, Err returns nil.
	// If Done is closed, Err returns a non-nil error explaining why:
	// Canceled if the context was canceled
	// or DeadlineExceeded if the context's deadline passed.
	// After Err returns a non-nil error, successive calls to Err return the same error.
	Err() error

	// Value returns the value associated with this context for key, or nil
	// if no value is associated with key. Successive calls to Value with
	// the same key returns the same result.
	//
	// Use context values only for request-scoped data that transits
	// processes and API boundaries, not for passing optional parameters to
	// functions.
	//
	// A key identifies a specific value in a Context. Functions that wish
	// to store values in Context typically allocate a key in a global
	// variable then use that key as the argument to context.WithValue and
	// Context.Value. A key can be any type that supports equality;
	// packages should define keys as an unexported type to avoid
	// collisions.
	//
	// Packages that define a Context key should provide type-safe accessors
	// for the values stored using that key:
	//
	// 	// Package user defines a User type that's stored in Contexts.
	// 	package user
	//
	// 	import "context"
	//
	// 	// User is the type of value stored in the Contexts.
	// 	type User struct {...}
	//
	// 	// key is an unexported type for keys defined in this package.
	// 	// This prevents collisions with keys defined in other packages.
	// 	type key int
	//
	// 	// userKey is the key for user.User values in Contexts. It is
	// 	// unexported; clients use user.NewContext and user.FromContext
	// 	// instead of using this key directly.
	// 	var userKey key
	//
	// 	// NewContext returns a new Context that carries value u.
	// 	func NewContext(ctx context.Context, u *User) context.Context {
	// 		return context.WithValue(ctx, userKey, u)
	// 	}
	//
	// 	// FromContext returns the User value stored in ctx, if any.
	// 	func FromContext(ctx context.Context) (*User, bool) {
	// 		u, ok := ctx.Value(userKey).(*User)
	// 		return u, ok
	// 	}
	Value(key any) any
}
```

### [context.WithCancel](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/context/context.go#L232)
```go
// WithCancel returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the returned cancel function is called
// or when the parent context's Done channel is closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := newCancelCtx(parent)
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, Canceled) }
}

// newCancelCtx returns an initialized cancelCtx.
func newCancelCtx(parent Context) cancelCtx {
	return cancelCtx{Context: parent}
}

// goroutines counts the number of goroutines ever created; for testing.
var goroutines atomic.Int32

// propagateCancel arranges for child to be canceled when parent is.
func propagateCancel(parent Context, child canceler) {
	done := parent.Done()
	if done == nil {
		return // parent is never canceled 先判断父节点状态，如果父节点永远不会触发取消时，会直接退出
	}

	select {
	case <-done:
		// parent is already canceled 判断父节点是否已经取消
		child.cancel(false, parent.Err())
		return
	default:
	}

	if p, ok := parentCancelCtx(parent); ok {  // 如果是go原生context
		p.mu.Lock()
		if p.err != nil { // 如果父节点已经取消，子节点会马上取消
			// parent has already been canceled
			child.cancel(false, p.err)
		} else { // 如果没有被取消，child 会被加入 parent 的 children 列表中，等待 parent 释放取消信号
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else { // 如果是自定义的context
		goroutines.Add(1)
		go func() { // 建立一个协程同时监听父和子两个节点的取消状态
			select { 
			case <-parent.Done(): // 父节点取消时调用子节点取消
				child.cancel(false, parent.Err()) 
			case <-child.Done():
			}
		}()
	}
}
```

### [context.WithTimeout](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/context/context.go#L232)
```go
// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete:
//
//	func slowOperationWithTimeout(ctx context.Context) (Result, error) {
//		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
//		defer cancel()  // releases resources if slowOperation completes before timeout elapses
//		return slowOperation(ctx)
//	}
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}
// WithDeadline returns a copy of the parent context with the deadline adjusted
// to be no later than d. If the parent's deadline is already earlier than d,
// WithDeadline(parent, d) is semantically equivalent to parent. The returned
// context's Done channel is closed when the deadline expires, when the returned
// cancel function is called, or when the parent context's Done channel is
// closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
    if parent == nil { // 判断父节点状态
		panic("cannot create context from nil parent")
    }
    if cur, ok := parent.Deadline(); ok && cur.Before(d) {
        // The current deadline is already sooner than the new one.  判断父节点是否有取消时间，如果父节点有取消时间而且在更早，则使用父节点的取消状态就行
        return WithCancel(parent)
    }
    c := &timerCtx{ // 创建一个timerCtx
        cancelCtx: newCancelCtx(parent),
        deadline:  d,
    }
    propagateCancel(parent, c) // 调用上面的cancel方法
    dur := time.Until(d)
    if dur <= 0 { // 判断是否已经到时间，如果已经到了的话，直接返回已经取消的ctx
        c.cancel(true, DeadlineExceeded) // deadline has already passed
        return c, func() { c.cancel(false, Canceled) }
    }
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.err == nil {
        c.timer = time.AfterFunc(dur, func() { // 添加timer的cancel方法
            c.cancel(true, DeadlineExceeded)    
        })
    }
    return c, func() { c.cancel(true, Canceled) }
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
    c.cancelCtx.cancel(false, err)
    if removeFromParent {
        // Remove this timerCtx from its parent cancelCtx's children. 从该ctx的父节点中删除该(子)节点
        removeChild(c.cancelCtx.Context, c)
    }
    c.mu.Lock()
    if c.timer != nil {
        c.timer.Stop() 
        c.timer = nil // timerCtx在取消时会主动释放定时器
    }
    c.mu.Unlock()
}
```

### todo:[context.WithValue](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/context/context.go#L523)
```go
// WithValue returns a copy of parent in which the value associated with key is
// val.
//
// Use context Values only for request-scoped data that transits processes and
// APIs, not for passing optional parameters to functions.
//
// The provided key must be comparable and should not be of type
// string or any other built-in type to avoid collisions between
// packages using context. Users of WithValue should define their own
// types for keys. To avoid allocating when assigning to an
// interface{}, context keys often have concrete type
// struct{}. Alternatively, exported context key variables' static
// type should be a pointer or interface.
func WithValue(parent Context, key, val any) Context {
	if parent == nil { // 判断父节点
		panic("cannot create context from nil parent")
	}
	if key == nil { // 判断key
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val} // 返回新context
}

// A valueCtx carries a key-value pair. It implements Value for that key and
// delegates all other calls to the embedded Context.
type valueCtx struct {  // 也就是说，一个valueCtx只存了一对k v
	Context
	key, val any
}

func (c *valueCtx) Value(key any) any { // 查找key
    if c.key == key { // 如果当前valueCtx存的key是要查找的key就直接返回值
        return c.val
    }
    return value(c.Context, key) // 调用value方法查询父节点的key
}

// &cancelCtxKey is the key that a cancelCtx returns itself for.
var cancelCtxKey int

func value(c Context, key any) any {
    for { // 循环遍历当前子节点上面的所有父节点，知道查找不到返回nil
        switch ctx := c.(type) {
        case *valueCtx: // 父节点类型是valueCtx
            if key == ctx.key { // 父节点如果查到对应kv，返回
                return ctx.val
            }
            c = ctx.Context // 如果不是对应kv，修改c为更上一层父节点，再次循环查询
        case *cancelCtx:
            if key == &cancelCtxKey { // 如果是key是获取当前ctx的cancelCtx的key
                return c // 返回cancelCtx用于调用cancel方法
            }
            c = ctx.Context // 如果不是对应kv，修改c为更上一层父节点，再次循环查询
        case *timerCtx:
            if key == &cancelCtxKey { // 如果是key是获取当前ctx的cancelCtx的key
                return &ctx.cancelCtx // 返回cancelCtx用于调用cancel方法
            }
            c = ctx.Context // 如果不是对应kv，修改c为更上一层父节点，再次循环查询
        case *emptyCtx: // 没有更上一层父节点了，返回Nil
            return nil
        default: // 如果是用户自定义context
            return c.Value(key) // 返回用户自定义Value方法
        }
    }
}
```
