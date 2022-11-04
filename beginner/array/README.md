# 数组

***
ctrl c/v永不过时 ！！！

> > 别的不说：先点先学：[菜鸟教程](https://www.runoob.com/go/go-arrays.html)

emmm~,在项目中用的一般比较少，一般我都是用常亮里面用来限制长度

### 1. 编译器优化：
**[传送门](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/cmd/compile/internal/walk/complit.go#L553)**
1. 当元素数量小于或者等于 4 个时，会直接将数组中的元素放置在栈上；
2. 当元素数量大于 4 个时，会将数组中的元素放置到静态区并在运行时取出；

```
func anylit(n ir.Node, var_ ir.Node, init *ir.Nodes) {
    t := n.Type()
    switch n.Op() {
        default:
        base.Fatalf("anylit: not lit, op=%v node=%v", n.Op(), n)
		...
        case ir.OSTRUCTLIT, ir.OARRAYLIT:
            n := n.(*ir.CompLitExpr)
            if !t.IsStruct() && !t.IsArray() {
                base.Fatalf("anylit: not struct/array")
            }
    
            if isSimpleName(var_) && len(n.List) > 4 {
                // lay out static data
                vstat := readonlystaticname(t)
    
                ctxt := inInitFunction
                if n.Op() == ir.OARRAYLIT {
                    ctxt = inNonInitFunction
                }
                fixedlit(ctxt, initKindStatic, n, vstat, init)
    
                // copy static to var
                appendWalkStmt(init, ir.NewAssignStmt(base.Pos, var_, vstat))
    
                // add expressions to automatic
                fixedlit(inInitFunction, initKindDynamic, n, var_, init)
                break
            }
    
            var components int64
            if n.Op() == ir.OARRAYLIT {
                components = t.NumElem()
            } else {
                components = int64(t.NumFields())
            }
            // initialization of an array or struct with unspecified components (missing fields or arrays)
            if isSimpleName(var_) || int64(len(n.List)) < components {
                appendWalkStmt(init, ir.NewAssignStmt(base.Pos, var_, nil))
            }
    
            fixedlit(inInitFunction, initKindLocalCode, n, var_, init)
        ...
    }
}
```