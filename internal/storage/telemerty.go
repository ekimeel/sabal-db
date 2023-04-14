package storage

import "github.com/ekimeel/sabal-pb/pb"

type Telemetry[V int | float64 | float32] struct {
	Equip     *pb.Equip `json:"equip"`
	Point     *pb.Point `json:"point"`
	Timestamp int64     `json:"ts"`
	Value     V         `json:"val"`
}
