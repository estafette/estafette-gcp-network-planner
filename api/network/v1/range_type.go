package network

type RangeType string

const (
	RangeTypePrimary   RangeType = "primary"
	RangeTypeSecondary RangeType = "secondary"

	RangeTypeUnknown RangeType = ""
)
