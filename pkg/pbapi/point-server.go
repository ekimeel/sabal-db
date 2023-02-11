package pbapi

import (
	"context"
	"errors"
	"github.com/ekimeel/sabal-db/pb"
	"github.com/ekimeel/sabal-db/pkg/services"
)

type PointServer struct {
	pb.UnimplementedPointServiceServer
	service services.PointService
}

func NewGrpcPointServer(pointService services.PointService) (*PointServer, error) {
	server := &PointServer{
		service: pointService,
	}

	return server, nil
}

func (s *PointServer) Get(ctx context.Context, uuid *pb.PointUUID) (*pb.Point, error) {
	val, ok := s.service.Get(uuid)
	if ok == false {
		return nil, errors.New("not found")
	}
	return val, nil
}

func (s *PointServer) GetAll(ctx context.Context, request *pb.ListRequest) (*pb.PointList, error) {
	values, err := s.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &pb.PointList{Points: values}, nil
}

func (s *PointServer) Create(ctx context.Context, point *pb.Point) (*pb.Point, error) {
	return s.service.Create(point)
}

func (s *PointServer) Update(ctx context.Context, point *pb.Point) (*pb.Point, error) {
	return s.service.Update(point)
}
func (s *PointServer) Delete(ctx context.Context, uuid *pb.PointUUID) (*pb.DeleteResponse, error) {
	err := s.service.Delete(uuid)
	if err != nil {
		return &pb.DeleteResponse{Success: false}, err
	}
	return &pb.DeleteResponse{Success: true}, nil
}

func (s *PointServer) GetAllByEquipUUID(ctx context.Context, uuid *pb.EquipUUID) (*pb.PointList, error) {
	values, err := s.service.GetAllByEquipUUID(uuid)
	if err != nil {
		return nil, err
	}

	return &pb.PointList{Points: values}, nil
}
