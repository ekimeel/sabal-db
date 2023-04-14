package services

import (
	"github.com/ekimeel/sabal-pb/pb"
)

type PointService interface {
	Get(id *pb.PointId) (*pb.Point, bool)
	GetAll() ([]*pb.Point, error)
	GetOrCreate(point *pb.Point) error
	GetAllByEquip(id *pb.EquipId) ([]*pb.Point, error)
	Create(point *pb.Point) error
	Update(point *pb.Point) error
	Delete(id *pb.PointId) error
}
