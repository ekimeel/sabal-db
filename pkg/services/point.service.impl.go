package services

import (
	"fmt"
	"github.com/ekimeel/sabal-db/internal/data"
	"github.com/ekimeel/sabal-db/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PointServiceImpl struct {
	data data.Collection[*pb.Point]
}

func NewPointService() PointService {
	s := &PointServiceImpl{}
	s.data = data.NewCollection[*pb.Point](0)
	return s
}

func createKeyWithPoint(point *pb.Point) string {
	return fmt.Sprintf("%s:%s", point.EquipUuid, point.Name)
}

func createKeyWithPointUUID(point *pb.PointUUID) string {
	return fmt.Sprintf("%s:%s", point.EquipUuid, point.Name)
}

func (s *PointServiceImpl) Create(point *pb.Point) (*pb.Point, error) {

	point.DateCreated = timestamppb.Now()
	point.LastUpdated = timestamppb.Now()

	err := s.data.Create(createKeyWithPoint(point), point)
	if err != nil {
		return point, nil
	}

	return point, err
}

func (s *PointServiceImpl) Update(point *pb.Point) (*pb.Point, error) {
	point.LastUpdated = timestamppb.Now()

	err := s.data.Replace(createKeyWithPoint(point), point)
	if err != nil {
		return point, nil
	}

	return point, err
}

func (s *PointServiceImpl) Get(uuid *pb.PointUUID) (*pb.Point, bool) {
	return s.data.Get(createKeyWithPointUUID(uuid))
}

func (s *PointServiceImpl) GetAll() ([]*pb.Point, error) {
	return s.data.Values(), nil
}

func (s *PointServiceImpl) Delete(uuid *pb.PointUUID) error {
	return s.data.Remove(createKeyWithPointUUID(uuid))
}

func (s *PointServiceImpl) GetAllByEquipUUID(equipUUID *pb.EquipUUID) ([]*pb.Point, error) {
	result := s.data.Filter(func(s string, point *pb.Point) bool {
		return point.EquipUuid == equipUUID.Id
	})

	return result, nil
}
