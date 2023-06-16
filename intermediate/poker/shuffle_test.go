package poker

import "testing"

func TestPoker_Shuffle(t *testing.T) {
	p.StorageBaseCards(Cards)
	p.ShuffleRandom()
	t.Log(p.CardToNum(p.cards))
}

var EnvCard2RealCard = map[int64]int64{3: 3, 4: 4, 5: 5, 6: 6, 7: 7,
	8: 8, 9: 9, 10: 10, 11: 11, 12: 12,
	13: 13, 14: 14, 15: 17, 16: 20, 17: 30}

func cardToDouzero(data *Card) int64 {
	return EnvCard2RealCard[data.Value]
}

func TestPoker_ShuffleProbability(t *testing.T) {
	p.StorageBaseCards(Cards)
	var req = [][3]int64{
		{PairStraight, 100, 2},
	}
	for i := 0; i < 100000; i++ {
		p.ShuffleProbability(req)
		if len(p.cards) != 54 {
			t.Fatal("总长度不够")
		}
		//var tmp = make([]int64, 17)
		//for i2 := range p.cards[0:17] {
		//	tmp[i2] = cardToDouzero(p.cards[i2])
		//}
		//t.Logf("%+v", tmp)
		//for i2 := range p.cards[17:34] {
		//	tmp[i2] = cardToDouzero(p.cards[17+i2])
		//}
		//t.Logf("%+v", tmp)
		//for i2 := range p.cards[34:51] {
		//	tmp[i2] = cardToDouzero(p.cards[34+i2])
		//}
		//t.Logf("%+v", tmp)
		//tmp = make([]int64, 3)
		//for i2 := range p.cards[51:54] {
		//	tmp[i2] = cardToDouzero(p.cards[51+i2])
		//}
		//t.Logf("%+v", tmp)
		var m = make(map[int64]struct{})
		for i2 := range p.cards {
			if _, ok := m[p.CardToNum(p.cards[i2])[0]]; !ok {
				m[p.CardToNum(p.cards[i2])[0]] = struct{}{}
			} else {
				t.Fatal("出现重复牌组")
			}
		}
	}
	t.Log("success")

	//p.ShuffleProbability(req)
	//t.Log(len(p.cards))
	////t.Log(result)
	//a := p.cards[:17]
	//p.SortCards(a)
	//t.Log(p.CardToNum(a))
	//
	//a = p.cards[17:34]
	//p.SortCards(a)
	//t.Log(p.CardToNum(a))
	//
	//a = p.cards[34:51]
	//p.SortCards(a)
	//t.Log(p.CardToNum(a))
	//
	//a = p.cards[51:]
	//p.SortCards(a)
	//t.Log(p.CardToNum(a))
}
