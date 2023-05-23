# poker 牌组封装
***
ctrl c/v永不过时 ！！！

```text
├─cardCombos.go     获取可用牌组，用于游戏中获取可以接或者出的牌
├─cards.go          牌组封装，找牌，出牌，卡牌当前状态等
├─cards_type.go     牌型判断
├─feature.go        特征值的计算
├─laizi.go          设置牌组中的癞子
├─poker.go          牌组封装类
├─shuffle.go        洗牌，支持可以通过牌型出现概率来洗牌内容(避免发出的牌过于零散)
├─value_set.go      牌值计算方式，用于从手牌中找到可用牌型
```