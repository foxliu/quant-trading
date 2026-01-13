package strategy

type PositionIntent string

const (
	Long  PositionIntent = "LONG"
	Short PositionIntent = "SHORT"
	Flat  PositionIntent = "FLAT" // 可选，但非常推荐
)
