# 哈希表

***
ctrl c/v永不过时 ！！！

> 别的不说：先点先学：[菜鸟教程](https://www.runoob.com/go/go-map.html)

### 1. 结构
##### 1. 基础结构
```go
// Map contains Type fields specific to maps.
type Map struct {
    Key  *Type // Key type 存放key
    Elem *Type // Val (elem) type 存放value
    
    Bucket *Type // internal struct type representing a hash bucket
    Hmap   *Type // internal struct type representing the Hmap (map header object)
    Hiter  *Type // internal struct type representing hash iterator state
}
```
##### 2. 底层结构
[传送门](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/runtime/map.go#L116)
```go
// A header for a Go map.
type hmap struct {
    // Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
    // Make sure this stays in sync with the compiler's definition.
    count     int // # live cells == size of map.  Must be first (used by len() builtin)
    flags     uint8
    B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
    noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
    hash0     uint32 // hash seed

    buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
    oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
    nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

    extra *mapextra // optional fields
}
```
- count: 字段表征了 map 目前的元素数目
- flags: flag 字段标志 map 的状态, 如 map 当前正在被遍历或正在被写入
- B: B 是哈希桶数目以 2 为底的对数, 在 go map 中, 哈希桶的数目都是 2 的整数次幂(这样设计的好处是可以是用位运算来计算取余运算的值, 即 N mod M = N & (M-1))
- noverflow: 是溢出桶的数目, 这个数值不是恒定精确的, 当其 B>=16 时为近似值
- hash0: hash0是随机哈希种子, map创建时调用 fastrand 函数生成的随机数, 设置的目的是为了降低哈希冲突的概率
- buckets: 是指向当前哈希桶的指针
- oldbuckets: 是当桶扩容时指向旧桶的指针
- nevacuate: 是当桶进行调整时指示的搬迁进度, 小于此地址的 buckets 是以前搬迁完毕的哈希桶
- extra: 是表征溢出桶的变量

其中buckets指向的桶是bmap,每一个 runtime.bmap 都能存储 8 个键值对，当哈希表中存储的数据过多，单个桶已经装满时就会使用 extra.nextOverflow 中桶存储溢出的数据。
```go
// A bucket for a Go map.
type bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	tophash [bucketCnt]uint8
	// Followed by bucketCnt keys and then bucketCnt elems.
	// NOTE: packing all the keys together and then all the elems together makes the
	// code a bit more complicated than alternating key/elem/key/elem/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}
```
在 Go 语言源代码中的定义只包含一个简单的 tophash 字段，tophash 存储了键的哈希的高 8 位，通过比较不同键的哈希的高 8 位可以减少访问键值对次数以提高性能
bmap在go的源码中并没有显式的定义出来，是因为其中数据是需要啊在编译期才能确定。不过通过反射的偏移量可以大致确定其结构
```go
type bmap struct {
    topbits  [8]uint8
    keys     [8]keytype
    elems    [8]elemtype
    //pad      uintptr(新的 go 版本已经移除了该字段, 我未具体了解此处的 change detail, 之前设置该字段是为了在 nacl/amd64p32 上的内存对齐)
    overflow uintptr
}
```
- topbits: 键哈希值的高 8 位
- keys: 哈希桶中所有键
- elems: 存放了哈希桶中的所有值
- overflow: 存放了所指向的溢出桶的地址

随着哈希表存储的数据逐渐增多，我们会扩容哈希表或者使用额外的桶存储溢出的数据，不会让单个桶中的数据超过 8 个，不过溢出桶只是临时的解决方案，创建过多的溢出桶最终也会导致哈希的扩容。
```go
// mapextra holds fields that are not present on all maps.
type mapextra struct {
	// If both key and elem do not contain pointers and are inline, then we mark bucket
	// type as containing no pointers. This avoids scanning such maps.
	// However, bmap.overflow is a pointer. In order to keep overflow buckets
	// alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.
	// overflow and oldoverflow are only used if key and elem do not contain pointers.
	// overflow contains overflow buckets for hmap.buckets.
	// oldoverflow contains overflow buckets for hmap.oldbuckets.
	// The indirection allows to store a pointer to the slice in hiter.
	overflow    *[]*bmap
	oldoverflow *[]*bmap

	// nextOverflow holds a pointer to a free overflow bucket.
	nextOverflow *bmap
}
```
如果key和elem都不包含指针并且是内联的，那么我们将桶类型标记为不包含指针。
这样可以避免扫描这样的地图。然而,bmap。Overflow是一个指针。
为了保持溢出桶存活，我们在hmap.extra.overflow和hmap.extra. oloverflow中存储了指向所有溢出桶的指针。
Overflow和oloverflow仅在key和elem不包含指针时使用。
Overflow包含hmap.buckets的溢出桶。
Oldoverflow包含了hmap.oldbuckets的溢出桶。间接允许在hiter中存储一个指向slice的指针。

### 2. 初始化
[传送门](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/cmd/compile/internal/walk/complit.go#L418)
```go
func maplit(n *ir.CompLitExpr, m ir.Node, init *ir.Nodes) {
	// make the map var
	a := ir.NewCallExpr(base.Pos, ir.OMAKE, nil, nil)
	a.SetEsc(n.Esc())
	a.Args = []ir.Node{ir.TypeNode(n.Type()), ir.NewInt(int64(len(n.List)))}
	litas(m, a, init)

	entries := n.List

	// The order pass already removed any dynamic (runtime-computed) entries.
	// All remaining entries are static. Double-check that.
	for _, r := range entries {
		r := r.(*ir.KeyExpr)
		if !isStaticCompositeLiteral(r.Key) || !isStaticCompositeLiteral(r.Value) {
			base.Fatalf("maplit: entry is not a literal: %v", r)
		}
	}

	if len(entries) > 25 {
        // For a large number of entries, put them in an array and loop.
        
        // build types [count]Tindex and [count]Tvalue
        tk := types.NewArray(n.Type().Key(), int64(len(entries)))
        te := types.NewArray(n.Type().Elem(), int64(len(entries)))
        
        // TODO(#47904): mark tk and te NoAlg here once the
        // compiler/linker can handle NoAlg types correctly.
        
        types.CalcSize(tk)
        types.CalcSize(te)
		······
	}
    // For a small number of entries, just add them directly.
    
    // Build list of var[c] = expr.
    // Use temporaries so that mapassign1 can have addressable key, elem.
    // TODO(josharian): avoid map key temporaries for mapfast_* assignments with literal keys.
    tmpkey := typecheck.Temp(m.Type().Key())
    tmpelem := typecheck.Temp(m.Type().Elem())
    
    for _, r := range entries {
        r := r.(*ir.KeyExpr)
        index, elem := r.Key, r.Value
        
        ir.SetPos(index)
        appendWalkStmt(init, ir.NewAssignStmt(base.Pos, tmpkey, index))
        
        ir.SetPos(elem)
        appendWalkStmt(init, ir.NewAssignStmt(base.Pos, tmpelem, elem))
        
        ir.SetPos(tmpelem)
        var a ir.Node = ir.NewAssignStmt(base.Pos, ir.NewIndexExpr(base.Pos, m, tmpkey), tmpelem)
        a = typecheck.Stmt(a) // typechecker rewrites OINDEX to OINDEXMAP
        a = orderStmtInPlace(a, map[string][]*ir.Name{})
        appendWalkStmt(init, a)
    }
    
    appendWalkStmt(init, ir.NewUnaryExpr(base.Pos, ir.OVARKILL, tmpkey))
    appendWalkStmt(init, ir.NewUnaryExpr(base.Pos, ir.OVARKILL, tmpelem))
```
**当哈希表中元素数量大于25个时，编译器会创建两个数组分别存储键和值。否则会直接将所有键值对一次加入到哈希表中**

### 2. 写入
[传送门](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/runtime/map.go#L578)
```go
// Like mapaccess, but allocates a slot for the key if it is not present in the map.
func mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
    hash := t.hasher(key, uintptr(h.hash0))
    
    // Set hashWriting after calling t.hasher, since t.hasher may panic,
    // in which case we have not actually done a write.
    h.flags ^= hashWriting
    
    if h.buckets == nil {
    h.buckets = newobject(t.bucket) // newarray(t.bucket, 1)
    }
    ...
}
```
- 获取hash值，判断buckets是否存在，不存在则创建
```go
again:
    bucket := hash & bucketMask(h.B)
    if h.growing() {
        growWork(t, h, bucket)
    }
    b := (*bmap)(add(h.buckets, bucket*uintptr(t.bucketsize)))
    top := tophash(hash)
	
	var inserti *uint8
	var insertk unsafe.Pointer
	var elem unsafe.Pointer
bucketloop:
	for {
		for i := uintptr(0); i < bucketCnt; i++ {
			if b.tophash[i] != top {
				if isEmpty(b.tophash[i]) && inserti == nil {
					inserti = &b.tophash[i]
					insertk = add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
					elem = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
				}
				if b.tophash[i] == emptyRest {
					break bucketloop
				}
				continue
			}
			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
			if t.indirectkey() {
				k = *((*unsafe.Pointer)(k))
			}
			if !t.key.equal(key, k) {
				continue
			}
			// already have a mapping for key. Update it.
			if t.needkeyupdate() {
				typedmemmove(t.key, k, key)
			}
			elem = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
			goto done
		}
		ovf := b.overflow(t)
		if ovf == nil {
			break
		}
		b = ovf
	}
```
**b.tophash[i] == emptyRestx**先通过tophash判断是否存在，再 **!t.key.equal(key, k)** 判断具体key的位置进行查询优化
```go
	// Did not find mapping for key. Allocate new cell & add entry.

	// If we hit the max load factor or we have too many overflow buckets,
	// and we're not already in the middle of growing, start growing.
	if !h.growing() && (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B)) {
		hashGrow(t, h)
		goto again // Growing the table invalidates everything, so try again
	}

	if inserti == nil {
		// The current bucket and all the overflow buckets connected to it are full, allocate a new one.
		newb := h.newoverflow(t, b)
		inserti = &newb.tophash[0]
		insertk = add(unsafe.Pointer(newb), dataOffset)
		elem = add(insertk, bucketCnt*uintptr(t.keysize))
	}

	// store new key/elem at insert position
	if t.indirectkey() {
		kmem := newobject(t.key)
		*(*unsafe.Pointer)(insertk) = kmem
		insertk = kmem
	}
	if t.indirectelem() {
		vmem := newobject(t.elem)
		*(*unsafe.Pointer)(elem) = vmem
	}
	typedmemmove(t.key, insertk, key)
	*inserti = top
	h.count++
```
- **!h.growing() && (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B))** 判断是否扩容
- **if inserti == nil** 判断是否创建新桶
- **typedmemmove(t.key, insertk, key)** 写入指针并添加写屏障避免gc

### 2. todo:扩容
[传送门](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/runtime/map.go#L657)
```go
func hashGrow(t *maptype, h *hmap) {
	// If we've hit the load factor, get bigger.
	// Otherwise, there are too many overflow buckets,
	// so keep the same number of buckets and "grow" laterally.
	bigger := uint8(1)
	if !overLoadFactor(h.count+1, h.B) {
        // 如果达到条件 1，那么将B值加1，相当于是原来的2倍
        // 否则对应条件 2，进行等量扩容，所以 B 不变
		bigger = 0
		h.flags |= sameSizeGrow
	}
	oldbuckets := h.buckets
	// 申请新的buckets空间
	newbuckets, nextOverflow := makeBucketArray(t, h.B+bigger, nil)

	flags := h.flags &^ (iterator | oldIterator)
	if h.flags&iterator != 0 {
		flags |= oldIterator
	}
	// commit the grow (atomic wrt gc)
	h.B += bigger
	h.flags = flags
	h.oldbuckets = oldbuckets
	h.buckets = newbuckets
    // 搬迁进度为0
	h.nevacuate = 0
    // overflow buckets 数为0
	h.noverflow = 0

	if h.extra != nil && h.extra.overflow != nil {
		// Promote current overflow buckets to the old generation.
		if h.extra.oldoverflow != nil {
			throw("oldoverflow is not nil")
		}
		h.extra.oldoverflow = h.extra.overflow
		h.extra.overflow = nil
	}
	if nextOverflow != nil {
		if h.extra == nil {
			h.extra = new(mapextra)
		}
		h.extra.nextOverflow = nextOverflow
	}

	// the actual copying of the hash table data is done incrementally
	// by growWork() and evacuate().
}

// overLoadFactor reports whether count items placed in 1<<B buckets is over loadFactor.
func overLoadFactor(count int, B uint8) bool {
	return count > bucketCnt && uintptr(count) > loadFactorNum*(bucketShift(B)/loadFactorDen)
}

// tooManyOverflowBuckets reports whether noverflow buckets is too many for a map with 1<<B buckets.
// Note that most of these overflow buckets must be in sparse use;
// if use was dense, then we'd have already triggered regular map growth.
func tooManyOverflowBuckets(noverflow uint16, B uint8) bool {
	// If the threshold is too low, we do extraneous work.
	// If the threshold is too high, maps that grow and shrink can hold on to lots of unused memory.
	// "too many" means (approximately) as many overflow buckets as regular buckets.
	// See incrnoverflow for more details.
	if B > 15 {
		B = 15
	}
	// The compiler doesn't see here that B < 16; mask B to generate shorter shift code.
	return noverflow >= uint16(1)<<(B&15)
}

func growWork(t *maptype, h *hmap, bucket uintptr) {
    // 为了确认搬迁的 bucket 是我们正在使用的 bucket
    // 即如果当前key映射到老的bucket1，那么就搬迁该bucket1。
    evacuate(t, h, bucket&h.oldbucketmask())
    
    // 如果还未完成扩容工作，则再搬迁一个bucket。
    if h.growing() {
    evacuate(t, h, h.nevacuate)
    }
}

```
触发条件：
1. 装载因子超过6.5(overLoadFactor 每次装载达到最大都会+1)
2. 过多的溢出桶(tooManyOverflowBuckets 根据当前的装载因子，扩容的最大装载因子不超过15)

因为如果对map不断的增删会造成overflow的bucket数量增多，但是负载因子并不高。针对上面两种情况所以有两种扩容方式：
- **增量扩容** ：对 1，将 B + 1，新建一个buckets数组，新的buckets大小是原来的2倍，然后旧buckets数据搬迁到新的buckets。
- **等量扩容** : 针对 2，并不扩大容量，buckets数量维持不变，重新做一遍类似增量扩容的搬迁动作，把松散的键值对重新排列一次，以使bucket的使用率更高，进而保证更快的存取。

### 2. todo:访问
[接受一个参数传送门](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/runtime/map.go#L395)

[接受两个个参数传送门](https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/runtime/map.go#L456)
- 当接受一个参数时，会使用 runtime.mapaccess1()，该函数仅会返回一个指向目标值的指针；
- 当接受两个参数时，会使用 runtime.mapaccess2()，除了返回目标值之外，它还会返回一个用于表示当前键对应的值是否存在的 bool 值：