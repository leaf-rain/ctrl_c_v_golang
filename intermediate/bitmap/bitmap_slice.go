package bitmap

type BitMapSwitchForSlice struct {
	bitmap []uint64 // 位图数据
}

// NewBitMapSwitch 初始化位图开关
// numSwitches: 初始化需要指定位图大小，避免无限制扩容
func NewBitMapSwitch(numSwitches int) *BitMapSwitchForSlice {
	numBits := (numSwitches + 63) / 64 // 计算所需的位图大小
	bitmap := make([]uint64, numBits)
	return &BitMapSwitchForSlice{
		bitmap: bitmap,
	}
}

// TurnOn 打开指定位置的开关
func (b *BitMapSwitchForSlice) TurnOn(switchIndex int) {
	if switchIndex >= 0 && switchIndex < len(b.bitmap)*64 {
		wordIndex := switchIndex / 64
		bitIndex := switchIndex % 64
		b.bitmap[wordIndex] |= (1 << bitIndex)
	}
}

// TurnOff 关闭指定位置的开关
func (b *BitMapSwitchForSlice) TurnOff(switchIndex int) {
	if switchIndex >= 0 && switchIndex < len(b.bitmap)*64 {
		wordIndex := switchIndex / 64
		bitIndex := switchIndex % 64
		b.bitmap[wordIndex] &= ^(1 << bitIndex)
	}
}

// IsOn 检查指定位置的开关状态
func (b *BitMapSwitchForSlice) IsOn(switchIndex int) bool {
	if switchIndex >= 0 && switchIndex < len(b.bitmap)*64 {
		wordIndex := switchIndex / 64
		bitIndex := switchIndex % 64
		return (b.bitmap[wordIndex] & (1 << bitIndex)) != 0
	}
	return false
}
