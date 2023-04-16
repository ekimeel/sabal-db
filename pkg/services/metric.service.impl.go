package services

import (
	"fmt"
	"github.com/ekimeel/sabal-pb/pb"
	"github.com/ekimeel/tstorage/pkg/tstorage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"
)

var metricServiceInstance *MetricServiceImpl
var metricServiceOnce sync.Once

type MetricServiceImpl struct {
	storage tstorage.Storage
}

func GetMetricService() MetricService {
	metricServiceOnce.Do(func() {
		log.Tracef("creating [%s] once", "metricServiceInstance")

		metricServiceInstance = &MetricServiceImpl{}

		var err error
		metricServiceInstance.storage, err = tstorage.NewStorage(
			tstorage.WithTimestampPrecision(tstorage.Seconds),
			tstorage.WithRetention((24*7)*time.Hour),
			tstorage.WithPartitionDuration(24*time.Hour),
			tstorage.WithDataPath("./data/"),
			tstorage.WithLogLevel(log.InfoLevel),
			tstorage.WithWalRecovery(tstorage.SkipAnyCorruptedRecord),
		)

		if err != nil {
			panic(err)
		}
	})

	return metricServiceInstance
}

func (s *MetricServiceImpl) Write(metric *pb.Metric) (bool, error) {
	err := s.storage.InsertRows([]tstorage.Row{
		{
			Metric: metric.PointId,
			DataPoint: tstorage.DataPoint{
				Timestamp: metric.Timestamp.Seconds,
				Value:     metric.Value,
			},
		},
	})
	if err != nil {
		return false, err
	}

	return true, err
}

func (s *MetricServiceImpl) WriteList(metricList *pb.MetricList) (bool, error) {
	rows := make([]tstorage.Row, len(metricList.Metrics))

	for i, metric := range metricList.Metrics {
		if metric.Timestamp == nil {
			return false, fmt.Errorf("nil timestamp at %d", i)
		}
		rows[i] = tstorage.Row{
			Metric: metric.PointId,
			DataPoint: tstorage.DataPoint{
				Timestamp: metric.Timestamp.Seconds,
				Value:     metric.Value,
			},
		}
	}

	err := s.storage.InsertRows(rows)
	if err != nil {
		return false, err
	}

	return true, err

}

func (s *MetricServiceImpl) Poll(id *pb.PointId) (*pb.Metric, error) {
	m := s.storage.Poll(id.Id)
	if m == nil {
		return nil, nil
	}
	return &pb.Metric{
		PointId:   id.Id,
		Timestamp: &timestamppb.Timestamp{Seconds: m.Timestamp},
		Value:     m.Value,
	}, nil

}
func (s *MetricServiceImpl) Select(id uint32, from, to int64) ([]*pb.Metric, error) {
	data, err := s.storage.Select(id, from, to)
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Metric, len(data))

	for i := 0; i < len(data); i++ {
		result[i] = &pb.Metric{
			PointId:   id,
			Timestamp: &timestamppb.Timestamp{Seconds: data[i].Timestamp},
			Value:     data[i].Value,
		}
	}

	return result, nil
}

func (s *MetricServiceImpl) Flush() error {
	return s.storage.Close()
}
