# 切片

***
ctrl c/v永不过时 ！！！

> 别的不说：先点先学：[菜鸟教程](https://www.runoob.com/go/go-slice.html)

### 1. 创建时校验

**[传送门](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/cmd/compile/internal/typecheck/type.go#L136)**

```go
// tcSliceType typechecks an OTSLICE node.
func tcSliceType(n *ir.SliceType) ir.Node {
    n.Elem = typecheckNtype(n.Elem)
    if n.Elem.Type() == nil {
    return n
    }
    t := types.NewSlice(n.Elem.Type())
    n.SetOTYPE(t)
    types.CheckSize(t)
    return n
}
```

### 2. 扩容方案

**[传送门](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/runtime/slice.go#L166)**
> 这里算法有所改动，1.之前是小于1024扩容会翻倍，现在这个值变成256，2. 大于固定值后扩容算法由之前的1.25倍换成现在的(length*768)/4

go1.18:
```go
newcap := old.cap
doublecap := newcap + newcap
if cap > doublecap {
    newcap = cap
} else {
    const threshold = 256
    if old.cap < threshold {
        newcap = doublecap
    } else {
        // Check 0 < newcap to detect overflow
        // and prevent an infinite loop.
        for 0 < newcap && newcap < cap {
            // Transition from growing 2x for small slices
            // to growing 1.25x for large slices. This formula
            // gives a smooth-ish transition between the two.
            newcap += (newcap + 3*threshold) / 4
        }
        // Set newcap to the requested cap when
        // the newcap calculation overflowed.
        if newcap <= 0 {
            newcap = cap
        }
    }
}
```
go1.18之前:
```go
newcap := old.cap
doublecap := newcap + newcap
if cap > doublecap {
    newcap = cap
} else {
    if old.len < 1024 {
        newcap = doublecap
    } else {
        for 0 < newcap && newcap < cap {
            newcap += newcap / 4
        }
        if newcap <= 0 {
            newcap = cap
        }
    }
}
```

