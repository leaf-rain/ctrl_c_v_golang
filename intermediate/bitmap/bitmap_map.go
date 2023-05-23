package bitmap

type BitMapSwitchByMap struct {
	bitmap map[int]uint64 // 位图数据
}

// NewBitMapSwitchMap 初始化位图开关
func NewBitMapSwitchMap() *BitMapSwitchByMap {
	return &BitMapSwitchByMap{
		bitmap: make(map[int]uint64),
	}
}

// TurnOn 打开指定位置的开关
func (b *BitMapSwitchByMap) TurnOn(switchIndex int) {
	if switchIndex >= 0 {
		wordIndex := switchIndex / 64
		bitIndex := switchIndex % 64
		if _, ok := b.bitmap[wordIndex]; !ok {
			b.bitmap[wordIndex] = 0
		}
		b.bitmap[wordIndex] |= (1 << bitIndex)
	}
}

// TurnOff 关闭指定位置的开关
func (b *BitMapSwitchByMap) TurnOff(switchIndex int) {
	if switchIndex >= 0 {
		wordIndex := switchIndex / 64
		bitIndex := switchIndex % 64
		if _, ok := b.bitmap[wordIndex]; !ok {
			b.bitmap[wordIndex] = 0
		}
		b.bitmap[wordIndex] &= ^(1 << bitIndex)
	}
}

// IsOn 检查指定位置的开关状态
func (b *BitMapSwitchByMap) IsOn(switchIndex int) bool {
	if switchIndex >= 0 {
		wordIndex := switchIndex / 64
		bitIndex := switchIndex % 64
		if _, ok := b.bitmap[wordIndex]; !ok {
			return false
		}
		return (b.bitmap[wordIndex] & (1 << bitIndex)) != 0
	}
	return false
}
