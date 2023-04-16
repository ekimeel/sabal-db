package services

import (
	"github.com/ekimeel/sabal-pb/pb"
)

type MetricService interface {
	Write(metric *pb.Metric) (bool, error)
	WriteList(metricList *pb.MetricList) (bool, error)
	Poll(id *pb.PointId) (*pb.Metric, error)
	Select(id uint32, to, from int64) ([]*pb.Metric, error)
	Flush() error
}
