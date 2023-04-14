package storage

type Statistic[V int | float64 | float32] struct {
	Average float64 `json:"avg,omitempty"`
	Minimum V       `json:"min,omitempty"`
	Maximum V       `json:"max,omitempty"`
	Sum     V       `json:"sum,omitempty"`
	Count   int     `json:"count,omitempty"`
	StdDev  float64 `json:"stddev,omitempty"`
	Median  V       `json:"med,omitempty"`
}

func NewStatistic[V int | float64 | float32]() Statistic[V] {
	return Statistic[V]{}
}
