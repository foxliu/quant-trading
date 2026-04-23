package strategy

// FROZEN: V1

type PositionIntent string

const (
	IntentLong  PositionIntent = "LONG"  // 做多 / 加多
	IntentShort PositionIntent = "SHORT" // 做空 / 加空
	IntentFlat  PositionIntent = "FLAT"  // 平仓
)
