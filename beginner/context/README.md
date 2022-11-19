# context上下文

***
ctrl c/v永不过时 ！！！好用请给star。

### 应用场景
常用于在并发程序中，用于处理超时、取消操作或者一些一场情况。
```go
// e.g:
func main() {
    messages := make(chan int, 10)
    done := make(chan bool)
    
    defer close(messages)
    // consumer
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        for _ = range ticker.C {
            select {
                case <-done:
                    fmt.Println("child process interrupt...")
                    return
                default:
                    fmt.Printf("send message: %d\n", <-messages)
            }
        }
    }()
    
    // producer
    for i := 0; i < 10; i++ {
        messages <- i
    }
    time.Sleep(5 * time.Second)
    close(done)
    time.Sleep(1 * time.Second)
    fmt.Println("main process exit!")
}
```


### 取消信号的传递
1. WithCancel
```go
// propagateCancel arranges for child to be canceled when parent is.
func propagateCancel(parent Context, child canceler) {
	done := parent.Done()
	if done == nil {
		return // parent is never canceled
	}

	select {
	case <-done:
		// parent is already canceled
		child.cancel(false, parent.Err())
		return
	default:
	}

	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil {
			// parent has already been canceled
			child.cancel(false, p.err)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		goroutines.Add(1)
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}
```
