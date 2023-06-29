package poker

import (
	"testing"
)

func Test_HintCardCombo(t *testing.T) {
	//p.SetLaizi([]int64{6})
	var c1 = []int64{81, 91, 101, 111, 121}
	var a, _, _, _ = p.DecodeFeature(41040134)
	t.Log("--->", a)
	feature := p.GetCardsFeature(c1, a)
	t.Log(feature)

	c1 = []int64{91, 101, 111, 112, 113, 114}
	result := p.HintCardCombo(c1, feature)
	t.Log(result)
	// 24050031
	//var cards = []int64{
	//	31, 32, 41, 43, 51, 52,
	//}
	//result := p.GetCardsFeature(cards, 0)
	//t.Log(result)

	//p.SetLaizi([]int64{9, 10})
	//cards = []int64{31,
	//	33,
	//	42,
	//	92,
	//	93,
	//	94,
	//	102,
	//	124,
	//	152}
	//result = p.GetCardsFeature(cards, 0)
	//t.Log(result)

}

func TestPoker_GetMinBomb(t *testing.T) {
	p.SetLaizi([]int64{7})
	//var tzz int64 = 41040134
	cards := p.NumToCard([]int64{31, 32, 33, 71, 41, 42, 51, 52})
	a, b, c, d := p.isFourWithTwoPair(cards)
	//a, b, c, d := p.isBomb(cards)
	t.Log(p.EncodeFeature(a, int(b), c, d))
}

func TestPoker_GetMinSingleStraight(t *testing.T) {
	//p.SetLaizi([]int64{6, 7})
	var c1 = []int64{61, 62}
	var a, _, _, _ = p.DecodeFeature(0)
	feature := p.GetCardsFeature(c1, a)
	t.Log(feature)

	c1 = []int64{61, 62, 71, 72, 81, 82, 91, 92}
	var cards = p.NumToCard(c1)
	p.UnUse(cards)
	result := p.GetMinOnePair(cards, feature, true, true, true)
	if result != nil {
		t.Log(p.CardToNum(result.Cards))
		a, b, c, d := p.isOnePair(result.Cards)
		t.Log("->", a, b, c, d)
		t.Log(p.EncodeFeature(a, int(b), c, d))
	} else {
		t.Errorf("failed")
	}

}
