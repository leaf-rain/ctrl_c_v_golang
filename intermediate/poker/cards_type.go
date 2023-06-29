package poker

const (
	Single         int64 = 1 // 单张
	OnePair        int64 = 2 // 对子
	Trio           int64 = 3 // 三张
	TrioWithSingle int64 = 4 // 三带单
	TrioWithPair   int64 = 5 // 三张带对子

	FourWithTwoSingle int64 = 8 // 四带两张单
	FourWithTwoPair   int64 = 9 // 四带两对

	SingleStraight         int64 = 21 // 单顺
	PairStraight           int64 = 22 // 连对
	TrioStraight           int64 = 23 // 飞机（不带牌）
	TrioStraightWithSingle int64 = 24 // 飞机带单张
	TrioStraightWithPair   int64 = 25 // 飞机带对子

	Bomb     int64 = 41 // 炸弹
	JokePair int64 = 51 // 王炸
)

//返回值：牌型，牌数(用于顺子或者炸弹计算长度)，牌值，是否带有癞子

// isBomb 是否为炸弹
func (p *Poker) isBomb(cards []*Card) (int64, int64, int64, int64) {
	if len(cards) < 4 { // 最少4张
		return 0, 0, 0, 0
	}
	vs, laiziCount := cardsToValueSet(cards)
	var length = len(vs)
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	if laiziCount == int64(len(cards)) { // 全是癞子的炸弹需要判断，如果全是同类型的癞子则当成没有癞子的炸弹打出
		fix = FixAll
		for i := range cards {
			if cards[i].Value != cards[0].Value {
				fix = FixBlend
			}
		}
	}
	var tmp int
	for i := range vs { // 除去赖子只可能剩1种牌型
		if vs[i].value > two {
			return 0, 0, 0, 0
		}
		if !vs[i].isLaizi {
			tmp++
			if tmp > 1 {
				return 0, 0, 0, 0
			}
		}
	}
	if length == 1 {
		return Bomb, int64(len(cards)), vs[0].value, fix
	} else {
		p.SortCards(cards)
		return Bomb, int64(len(cards)), cards[len(cards)-1].Value, fix
	}
}

// isJokePair 是否为王炸
func (p *Poker) isJokePair(nCards []*Card) (int64, int64, int64, int64) {
	if len(nCards) < 2 {
		return 0, 0, 0, 0
	}
	var result int64
	var isContainBigKing bool
	for _, item := range nCards {
		if item.Value != littleKing && item.Value != bigKing {
			return 0, 0, 0, 0
		}
		if result < item.Value {
			result = item.Value
		}
		if item.Value == bigKing {
			isContainBigKing = true
			result++
		}
	}
	if isContainBigKing {
		result--
	}
	return JokePair, int64(len(nCards) - 1), result, FixNo
}

// isSingle 是否为单张
func (p *Poker) isSingle(nCards []*Card) (int64, int64, int64, int64) {
	if len(nCards) != 1 || nCards[0] == nil {
		return 0, 0, 0, 0
	}
	var laiziCount int64
	if nCards[0].IsLaizi {
		laiziCount++
	}
	var fix = FixNo
	if nCards[0].IsLaizi {
		fix = FixHave
	}
	return Single, 1, nCards[0].Value, fix
}

// isOnePair 是否为一对
func (p *Poker) isOnePair(cards []*Card) (int64, int64, int64, int64) {
	if len(cards) != 2 {
		return 0, 0, 0, 0
	}
	vs, laiziCount := cardsToValueSet(cards)
	if len(vs) > 2 { // 除去赖子只可能剩两种牌型
		return 0, 0, 0, 0
	}
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	// 判断癞子数量
	if laiziCount == 2 {
		if tmpValue := p.GetMaxNoJoker(cards); tmpValue != 0 {
			return OnePair, 1, tmpValue, fix
		}
		return OnePair, 1, cards[1].Value, fix
	}
	p.SortCards(cards)
	if laiziCount == 1 && vs[0].value < littleKing {
		return OnePair, 1, vs[0].value, fix
	} else {
		// 没有癞子
		if cards[0].Value == cards[1].Value &&
			cards[0].Value < littleKing {
			return OnePair, 1, cards[0].Value, fix
		}
	}
	return 0, 0, 0, 0
}

// isTrio 是否为三条
func (p *Poker) isTrio(cards []*Card) (int64, int64, int64, int64) {
	if len(cards) != 3 {
		return 0, 0, 0, 0
	}
	vs, laiziCount := cardsToValueSet(cards)
	if len(vs) > 2 { // 除去赖子只可能剩两种牌型
		return 0, 0, 0, 0
	}
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	// 判断癞子数量
	if laiziCount == 3 { // 3张全是癞子牌
		if tmpValue := p.GetMaxNoJoker(cards); tmpValue != 0 {
			return Trio, 1, tmpValue, fix
		}
		return Trio, 1, cards[2].Value, fix
	}
	p.SortCards(cards)
	if laiziCount == 2 { // 有癞子牌
		var maxValue int64
		for i := range vs {
			if vs[i].isLaizi {
				continue
			}
			if vs[i].value > maxValue && vs[i].value < littleKing {
				maxValue = vs[i].value
			}
		}
		return Trio, 1, maxValue, fix
	}
	if laiziCount == 1 { // 有癞子牌
		var tmpValue int
		var maxValue int64
		for i := range vs {
			if vs[i].isLaizi {
				continue
			}
			tmpValue++
			if vs[i].value > maxValue && vs[i].value < littleKing {
				maxValue = vs[i].value
			}
		}
		if tmpValue > 1 {
			return 0, 0, 0, 0
		} else {
			return Trio, 1, maxValue, fix
		}
	} else {
		// 没有癞子 判断三张是否都相等
		if cards[0].Value == cards[1].Value &&
			cards[1].Value == cards[2].Value {
			return Trio, 1, cards[0].Value, fix
		}
	}
	return 0, 0, 0, 0
}

// isTrioWithSingle 是否为三带一
func (p *Poker) isTrioWithSingle(cards []*Card) (int64, int64, int64, int64) {
	if len(cards) != 4 {
		return 0, 0, 0, 0
	}
	vs, laiziCount := cardsToValueSet(cards)
	if len(vs) > 2 { // 除去赖子只可能剩两种牌型
		return 0, 0, 0, 0
	}
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	p.SortCards(cards)
	if laiziCount == 4 {
		if tmpValue := p.GetMaxNoJoker(cards); tmpValue != 0 {
			return TrioWithSingle, 1, tmpValue, fix
		}
		return TrioWithSingle, 1, cards[3].Value, fix
	} else if laiziCount == 3 && vs[0].value < littleKing {
		return TrioWithSingle, 1, vs[0].value, fix
	} else if laiziCount == 2 {
		if len(vs) == 1 && vs[0].value < littleKing {
			return TrioWithSingle, 1, vs[0].value, fix
		} else {
			var result int64
			for _, item := range vs {
				if result < item.value && item.value < littleKing {
					result = item.value
				}
			}
			if result > 0 {
				return TrioWithSingle, 1, result, fix
			}
		}
	} else if laiziCount == 1 {
		if tmpVs := getValueSetByTimesNoJoker(vs); tmpVs != nil {
			return TrioWithSingle, 1, tmpVs.value, fix
		}
	} else {
		if tmpVs := getValueSetByTimesNoJoker(vs); tmpVs != nil && tmpVs.times >= 3 {
			return TrioWithSingle, 1, tmpVs.value, fix
		}
	}
	return 0, 0, 0, 0
}

// isTrioWithPair 是否为三带对
func (p *Poker) isTrioWithPair(cards []*Card) (int64, int64, int64, int64) {
	if len(cards) != 5 {
		return 0, 0, 0, 0
	}
	vs, laiziCount := cardsToValueSet(cards)
	if len(vs) > 2 { // 除去赖子只可能剩两种牌型
		return 0, 0, 0, 0
	}
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	p.SortCards(cards)
	if laiziCount >= 3 { // 3张癞子以上，直接拿最大的非王值
		var result = p.GetMaxNoJoker(cards)
		if result > 0 {
			return TrioWithPair, 1, result, fix
		}
	} else if laiziCount >= 1 {
		if vs[0].times == 2 { // 两个对子，选大的对子
			if tmpVs := getValueSetByValueNoJoker(vs); tmpVs != nil {
				return TrioWithPair, 1, tmpVs.value, fix
			}
		} else {
			if tmpVs := getValueSetByTimesNoJoker(vs); tmpVs != nil {
				return TrioWithPair, 1, tmpVs.value, fix
			}
		}
	} else {
		if tmpVs := getValueSetByTimesNoJoker(vs); tmpVs != nil && tmpVs.times == 3 {
			return TrioWithPair, 1, tmpVs.value, fix
		}
	}
	return 0, 0, 0, 0
}

// isFourWithTwoSingle 判断是否为四带两张单牌
func (p *Poker) isFourWithTwoSingle(cards []*Card) (int64, int64, int64, int64) {
	if len(cards) != 6 {
		return 0, 0, 0, 0
	}
	vs, laiziCount, noLaizi := cardsToValueSetOnLaizi(cards)
	if noLaizi > 3 { // 除去赖子只可能剩3种牌型
		return 0, 0, 0, 0
	}
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	if int64(len(cards)) == laiziCount { // 全是癞子牌
		var result = p.GetMaxNoJoker(cards)
		if result > 0 {
			return FourWithTwoSingle, 1, result, fix
		}
	}
	// 用值大小来排序
	valueSetSortByValue(vs)
	var tmpLaiziCount int64
	for _, item := range vs {
		tmpLaiziCount = laiziCount
		if item.isLaizi {
			tmpLaiziCount -= item.times
		}
		if item.value >= littleKing && item.times > 2 { // 王炸不能带牌
			return 0, 0, 0, 0
		}
		if item.value >= littleKing {
			continue
		}
		if item.times+tmpLaiziCount >= 4 {
			return FourWithTwoSingle, 1, item.value, fix
		}
	}
	return 0, 0, 0, 0
}

// isFourWithTwoPair 是否为四张带两对
func (p *Poker) isFourWithTwoPair(cards []*Card) (int64, int64, int64, int64) {
	if len(cards) != 8 {
		return 0, 0, 0, 0
	}
	vs, laiziCount, noLaizi := cardsToValueSetOnLaizi(cards)
	if noLaizi > 3 { // 除去赖子只可能剩3种牌型
		return 0, 0, 0, 0
	}
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	if int64(len(cards)) == laiziCount { // 全是癞子牌
		p.SortCards(cards)
		return FourWithTwoPair, 1, cards[len(cards)-1].Value, fix
	}
	// 获取最大的值
	var result int64
	var flag bool
	// 用值大小来排序
	valueSetSortByValue(vs)
	var tmpLaiziCount int64
	for _, item := range vs {
		tmpLaiziCount = laiziCount
		if item.value >= littleKing {
			continue
		}
		if item.isLaizi {
			tmpLaiziCount -= item.times
		}
		if item.times < 4 {
			tmpLaiziCount -= 4 - item.times
		}
		if result < item.value && tmpLaiziCount >= 0 {
			flag = true
			for _, item2 := range vs {
				if item2.value == item.value || item2.isLaizi {
					continue
				}
				if item2.times%2 == 0 {
					continue
				} else {
					tmpLaiziCount--
					if tmpLaiziCount < 0 {
						flag = false
						break
					}
				}
			}
			if flag {
				result = item.value
			}
		}
	}
	if result > 0 {
		return FourWithTwoPair, 1, result, fix
	} else {
		return 0, 0, 0, 0
	}
}

// isSingleStraight 是否为单张顺子
func (p *Poker) isSingleStraight(cards []*Card) (int64, int64, int64, int64) {
	if len(cards) < 5 {
		return 0, 0, 0, 0
	}
	var section = int64(len(cards))
	vs, laiziCount := cardsToValueMap(cards)
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	if int64(len(cards)) == laiziCount { // 全是癞子牌
		p.SortCards(cards)
		return SingleStraight, 1, ace, fix
	}
	var length = len(cards)
	var result int64
	var ok bool
	var timLaiziCount int64
	var flag bool
	for i := ace; i >= 0; i-- {
		timLaiziCount = laiziCount
		result = i
		_, ok = vs[i]
		if !ok && timLaiziCount < 1 {
			continue
		}
		if !ok || vs[i].isLaizi {
			timLaiziCount--
		}
		for i2 := 1; i2 < length; i2++ {
			if _, ok = vs[result-int64(i2)]; !ok {
				timLaiziCount--
				if timLaiziCount < 0 {
					break
				}
			} else {
				if vs[result-int64(i2)].isLaizi {
					timLaiziCount--
					if timLaiziCount < 0 {
						break
					}
				}
			}
			if i2 == length-1 {
				flag = true
				break
			}
		}
		if flag {
			break
		}
	}
	if flag {
		return SingleStraight, section, result, fix
	}
	return 0, 0, 0, 0
}

// isPairStraight 是否为双顺
func (p *Poker) isPairStraight(cards []*Card) (int64, int64, int64, int64) {
	var length = len(cards)
	if length < 6 || length%2 != 0 {
		return 0, 0, 0, 0
	}
	var section = int64(len(cards)) / 2
	vs, laiziCount := cardsToValueMap(cards)
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	if int64(len(cards)) == laiziCount { // 全是癞子牌
		p.SortCards(cards)
		return PairStraight, 1, ace, fix
	}
	length = len(cards) / 2
	var result int64
	var ok bool
	var timLaiziCount int64
	var flag bool
	for i := ace; i >= 0; i-- {
		timLaiziCount = laiziCount
		result = i
		_, ok = vs[i]
		if !ok && timLaiziCount < 2 {
			continue
		}
		if !ok || vs[i].isLaizi {
			timLaiziCount -= 2
		}
		if ok && vs[result].times < 2 {
			timLaiziCount -= (2 - vs[result].times)
			if timLaiziCount < 0 {
				continue
			}
		}
		for i2 := 1; i2 < length; i2++ {
			if _, ok = vs[result-int64(i2)]; !ok {
				timLaiziCount -= 2
				if timLaiziCount < 0 {
					break
				}
			} else {
				if vs[result-int64(i2)].isLaizi {
					timLaiziCount -= 2
					if timLaiziCount < 0 {
						break
					}
				} else {
					if vs[result-int64(i2)].times < 2 {
						timLaiziCount -= (2 - vs[result-int64(i2)].times)
						if timLaiziCount < 0 {
							break
						}
					}
				}
			}
			if i2 == length-1 {
				flag = true
				break
			}
		}
		if flag {
			break
		}
	}
	if flag {
		return PairStraight, section, result, fix
	}
	return 0, 0, 0, 0
}

// isTrioStraight 是否为三顺
func (p *Poker) isTrioStraight(cards []*Card) (int64, int64, int64, int64) {
	var length = len(cards)
	if length < 6 || length%3 != 0 {
		return 0, 0, 0, 0
	}
	var section = int64(len(cards)) / 3
	vs, laiziCount := cardsToValueMap(cards)
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	if int64(len(cards)) == laiziCount { // 全是癞子牌
		p.SortCards(cards)
		return TrioStraight, 1, ace, fix
	}
	var l int64 = 3
	length = len(cards) / int(l)
	var result int64
	var ok bool
	var timLaiziCount int64
	var flag bool
	for i := ace; i >= 0; i-- {
		timLaiziCount = laiziCount
		result = i
		_, ok = vs[i]
		if !ok && timLaiziCount < l {
			continue
		}
		if !ok || vs[i].isLaizi {
			timLaiziCount -= l
		}
		if ok && vs[result].times < l {
			timLaiziCount -= (l - vs[result].times)
			if timLaiziCount < 0 {
				continue
			}
		}
		for i2 := 1; i2 < length; i2++ {
			if _, ok = vs[result-int64(i2)]; !ok {
				timLaiziCount -= l
				if timLaiziCount < 0 {
					break
				}
			} else {
				if vs[result-int64(i2)].isLaizi {
					timLaiziCount -= l
					if timLaiziCount < 0 {
						break
					}
				} else {
					if vs[result-int64(i2)].times < l {
						timLaiziCount -= (l - vs[result-int64(i2)].times)
						if timLaiziCount < 0 {
							break
						}
					}
				}
			}
			if i2 == length-1 {
				flag = true
				break
			}
		}
		if flag {
			break
		}
	}
	if flag {
		return TrioStraight, section, result, fix
	}
	return 0, 0, 0, 0
}

// isTrioStraightWithSingle 是否为飞机带单牌
func (p *Poker) isTrioStraightWithSingle(cards []*Card) (int64, int64, int64, int64) {
	var length = len(cards)
	if length < 8 || length%4 != 0 {
		return 0, 0, 0, 0
	}
	var section = int64(len(cards)) / 4
	vs, laiziCount := cardsToValueMap(cards)
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	length = len(cards) / 4
	if laiziCount >= int64(len(cards)-length) { // 全是癞子牌
		p.SortCards(cards)
		return TrioStraightWithSingle, 1, ace, fix
	}
	var l int64 = 3
	var result int64
	var ok bool
	var timLaiziCount int64
	var flag bool
	for i := ace; i >= 0; i-- {
		timLaiziCount = laiziCount
		result = i
		_, ok = vs[i]
		if !ok && timLaiziCount < l {
			continue
		}
		if !ok || vs[i].isLaizi {
			timLaiziCount -= l
		}
		if ok && vs[result].times < l {
			timLaiziCount -= (l - vs[result].times)
			if timLaiziCount < 0 {
				continue
			}
		}
		for i2 := 1; i2 < length; i2++ {
			if _, ok = vs[result-int64(i2)]; !ok {
				timLaiziCount -= l
				if timLaiziCount < 0 {
					break
				}
			} else {
				if vs[result-int64(i2)].isLaizi {
					timLaiziCount -= l
					if timLaiziCount < 0 {
						break
					}
				} else {
					if vs[result-int64(i2)].times < l {
						timLaiziCount -= (l - vs[result-int64(i2)].times)
						if timLaiziCount < 0 {
							break
						}
					}
				}
			}
			if i2 == length-1 {
				flag = true
				break
			}
		}
		if flag {
			break
		}
	}
	if flag {
		return TrioStraightWithSingle, section, result, fix
	}
	return 0, 0, 0, 0
}

// isTrioStraightWithPair 是否为飞机带对子
func (p *Poker) isTrioStraightWithPair(cards []*Card) (int64, int64, int64, int64) {
	var length = len(cards)
	if length < 10 || length%5 != 0 {
		return 0, 0, 0, 0
	}
	var section = int64(len(cards)) / 5
	vs, laiziCount := cardsToValueMap(cards)
	var fix = FixNo
	if laiziCount > 0 {
		fix = FixHave
	}
	length = len(cards) / 5
	if laiziCount >= int64(len(cards)-length) { // 癞子牌过多不同判断
		p.SortCards(cards)
		return TrioStraightWithPair, 1, ace, fix
	}
	var l int64 = 3
	var result int64
	var ok bool
	var timLaiziCount int64
	var flag bool
	for i := ace; i >= 0; i-- {
		timLaiziCount = laiziCount
		result = i
		_, ok = vs[i]
		if !ok && timLaiziCount < l {
			continue
		}
		var outed = make(map[int64]int64)
		if !ok || vs[i].isLaizi {
			timLaiziCount -= l
		}
		if ok && vs[result].times < l {
			timLaiziCount -= (l - vs[result].times)
			if timLaiziCount < 0 {
				continue
			}
			outed[vs[i].value] = vs[result].times
		} else if ok && vs[result].times >= l {
			outed[vs[i].value] = 3
		}
		var trioFlag bool
		for i2 := 1; i2 < length; i2++ {
			if _, ok = vs[result-int64(i2)]; !ok {
				timLaiziCount -= l
				if timLaiziCount < 0 {
					break
				}
			} else {
				if vs[result-int64(i2)].isLaizi {
					timLaiziCount -= l
					if timLaiziCount < 0 {
						break
					}
				} else {
					if vs[result-int64(i2)].times < l {
						timLaiziCount -= (l - vs[result-int64(i2)].times)
						if timLaiziCount < 0 {
							break
						}
					}
				}
				if vs[result-int64(i2)].times >= l {
					outed[vs[result-int64(i2)].value] = 3
				} else {
					outed[vs[result-int64(i2)].value] = vs[result-int64(i2)].times
				}
			}
			if i2 == length-1 {
				trioFlag = true
				break
			}
		}
		if trioFlag { // 判断剩下的牌是否全是对子
			var pair int // 对子数量
			var tmpTimes int64
			for k, v := range vs {
				if v.isLaizi {
					continue
				}
				if _, ok = outed[k]; !ok {
					tmpTimes = v.times
				} else {
					tmpTimes = v.times - outed[k]
				}
				if tmpTimes <= 0 {
					continue
				}
				if tmpTimes%2 != 0 {
					timLaiziCount -= 1
					if timLaiziCount < 0 {
						break
					}
					tmpTimes += 1
				}
				pair += int(tmpTimes / 2)
			}
			pair += int(timLaiziCount / 2)
			if pair == length {
				flag = true
				break
			}
		}
	}
	if flag {
		return TrioStraightWithPair, section, result, fix
	}
	return 0, 0, 0, 0
}
