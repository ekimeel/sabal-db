package pbapi

import (
	"context"
	"errors"
	"github.com/ekimeel/sabal-db/pkg/services"
	"github.com/ekimeel/sabal-pb/pb"
)

type EquipServer struct {
	pb.UnimplementedEquipServiceServer
	service services.EquipService
}

func NewGrpcEquipServer(equipService services.EquipService) (*EquipServer, error) {
	server := &EquipServer{
		service: equipService,
	}

	return server, nil
}

func (s *EquipServer) Get(ctx context.Context, uuid *pb.EquipId) (*pb.Equip, error) {
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

func (s *EquipServer) Create(ctx context.Context, equip *pb.Equip) (*pb.Equip, error) {
	err := s.service.Create(equip)
	return equip, err
}

func (s *EquipServer) Update(ctx context.Context, equip *pb.Equip) (*pb.Equip, error) {
	err := s.service.Update(equip)
	return equip, err
}
func (s *EquipServer) Delete(ctx context.Context, uuid *pb.EquipId) (*pb.DeleteResponse, error) {
	err := s.service.Delete(uuid)
	if err != nil {
		return &pb.DeleteResponse{Success: false}, err
	}
	return &pb.DeleteResponse{Success: true}, nil
}

func (s *EquipServer) GetOrCreate(ctx context.Context, equip *pb.Equip) (*pb.Equip, error) {
	err := s.service.GetOrCreate(equip)
	return equip, err
}
