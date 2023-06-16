package poker

import (
	"testing"
)

func TestPoker_CompareFeature(t *testing.T) {
	t.Log(p.CompareFeature(1010134, 1010044))
}

func TestPoker_TestFeature(t *testing.T) {
	p.SetLaizi([]int64{4})
	var cd = []int64{41, 42}
	var cs = p.NumToCard(cd)
	a, b, c, d := p.isOnePair(cs)
	t.Log(p.EncodeFeature(a, int(b), c, d))
}

func TestPoker_GetCardsFeature(t *testing.T) {
	//p.SetLaizi([]int64{6, 7})
	var cards = []int64{
		31, 32, 33, 41, 42, 43, 111, 112, 121, 122,
	}
	var result = p.GetCardsFeature(cards, 0)
	t.Log(result)
}
