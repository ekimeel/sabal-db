package services

import "github.com/ekimeel/sabal-db/pb"

type EquipService interface {
	Get(uuid *pb.EquipUUID) (*pb.Equip, bool)
	GetAll() ([]*pb.Equip, error)
	Create(point *pb.Equip) (*pb.Equip, error)
	Update(point *pb.Equip) (*pb.Equip, error)
	Delete(uuid *pb.EquipUUID) error
}
