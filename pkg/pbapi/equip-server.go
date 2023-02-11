package pbapi

import (
	"context"
	"errors"
	"github.com/ekimeel/sabal-db/pb"
	"github.com/ekimeel/sabal-db/pkg/services"
)

type EquipServer struct {
	pb.UnimplementedEquipServiceServer
	service services.EquipService
}

func NewGrpcEquipServer(thingService services.EquipService) (*EquipServer, error) {
	server := &EquipServer{
		service: thingService,
	}

	return server, nil
}

func (s *EquipServer) Get(ctx context.Context, uuid *pb.EquipUUID) (*pb.Equip, error) {
	val, ok := s.service.Get(uuid)
	if ok == false {
		return nil, errors.New("not found")
	}
	return val, nil
}

func (s *EquipServer) GetAll(ctx context.Context, request *pb.ListRequest) (*pb.EquipList, error) {
	values, err := s.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &pb.EquipList{Equips: values}, nil
}

func (s *EquipServer) Create(ctx context.Context, point *pb.Equip) (*pb.Equip, error) {
	return s.service.Create(point)
}

func (s *EquipServer) Update(ctx context.Context, point *pb.Equip) (*pb.Equip, error) {
	return s.service.Update(point)
}
func (s *EquipServer) Delete(ctx context.Context, uuid *pb.EquipUUID) (*pb.DeleteResponse, error) {
	err := s.service.Delete(uuid)
	if err != nil {
		return &pb.DeleteResponse{Success: false}, err
	}
	return &pb.DeleteResponse{Success: true}, nil
}
