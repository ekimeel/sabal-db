package services

import (
	"github.com/ekimeel/sabal-db/pb"
)

type PointService interface {
	Get(uuid *pb.PointUUID) (*pb.Point, bool)
	GetAll() ([]*pb.Point, error)
	Create(point *pb.Point) (*pb.Point, error)
	Update(point *pb.Point) (*pb.Point, error)
	Delete(uuid *pb.PointUUID) error
	GetAllByEquipUUID(uuid *pb.EquipUUID) ([]*pb.Point, error)
}
