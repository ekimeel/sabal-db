package pbapi

import (
	"context"
	"errors"
	"github.com/ekimeel/sabal-db/pkg/services"
	"github.com/ekimeel/sabal-pb/pb"
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

func (s *PointServer) Get(ctx context.Context, id *pb.PointId) (*pb.Point, error) {
	val, ok := s.service.Get(id)
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
	err := s.service.Create(point)
	return point, err
}

func (s *PointServer) Update(ctx context.Context, point *pb.Point) (*pb.Point, error) {
	err := s.service.Update(point)
	return point, err
}
func (s *PointServer) Delete(ctx context.Context, id *pb.PointId) (*pb.DeleteResponse, error) {
	err := s.service.Delete(id)
	if err != nil {
		return &pb.DeleteResponse{Success: false}, err
	}
	return &pb.DeleteResponse{Success: true}, nil
}

func (s *PointServer) GetAllByEquip(ctx context.Context, id *pb.EquipId) (*pb.PointList, error) {
	values, err := s.service.GetAllByEquip(id)
	if err != nil {
		return nil, err
	}

	return &pb.PointList{Points: values}, nil
}

func (s *PointServer) GetOrCreate(ctx context.Context, point *pb.Point) (*pb.Point, error) {
	err := s.service.GetOrCreate(point)
	return point, err
}
