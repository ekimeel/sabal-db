package services

import (
	"github.com/ekimeel/sabal-db/internal/data"
	"github.com/ekimeel/sabal-db/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EquipServiceImpl struct {
	data data.Collection[*pb.Equip]
}

func NewEquipService() EquipService {
	s := &EquipServiceImpl{}
	s.data = data.NewCollection[*pb.Equip](0)
	return s
}

func (s *EquipServiceImpl) Create(point *pb.Equip) (*pb.Equip, error) {

	point.DateCreated = timestamppb.Now()
	point.LastUpdated = timestamppb.Now()

	err := s.data.Create(point.Uuid, point)
	if err != nil {
		return point, nil
	}

	return point, err
}

func (s *EquipServiceImpl) Update(point *pb.Equip) (*pb.Equip, error) {
	point.LastUpdated = timestamppb.Now()

	err := s.data.Replace(point.Uuid, point)
	if err != nil {
		return point, nil
	}

	return point, err
}

func (s *EquipServiceImpl) Get(uuid *pb.EquipUUID) (*pb.Equip, bool) {
	return s.data.Get(uuid.GetId())
}

func (s *EquipServiceImpl) GetAll() ([]*pb.Equip, error) {
	return s.data.Values(), nil
}

func (s *EquipServiceImpl) Delete(uuid *pb.EquipUUID) error {
	return s.data.Remove(uuid.GetId())
}
