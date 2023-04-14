package pbapi

import (
	"context"
	"github.com/ekimeel/sabal-db/pkg/services"
	"github.com/ekimeel/sabal-pb/pb"
)

type MetricServer struct {
	pb.UnimplementedMetricServiceServer
	service services.MetricService
}

func (s *MetricServer) Write(ctx context.Context, metric *pb.Metric) (*pb.MetricWriteResponse, error) {
	r := &pb.MetricWriteResponse{Accepted: false}
	accepted, err := s.service.Write(metric)
	if err != nil {
		return r, err
	}
	r.Accepted = accepted

	return r, nil
}

func (s *MetricServer) Poll(ctx context.Context, id *pb.PointId) (*pb.Metric, error) {
	return s.service.Poll(id)
}

func (s *MetricServer) Select(ctx context.Context, request *pb.MetricRequest) (*pb.MetricList, error) {

	items, err := s.service.Select(request.PointId, request.From.Seconds, request.To.Seconds)
	if err != nil {
		return nil, err
	}

	resp := &pb.MetricList{}
	resp.Metrics = items
	return resp, nil
}

func NewGrpcMetricServer(metricService services.MetricService) (*MetricServer, error) {
	server := &MetricServer{
		service: metricService,
	}

	return server, nil
}
