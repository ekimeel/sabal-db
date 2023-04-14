package services

import "github.com/ekimeel/sabal-pb/pb"

type EquipService interface {
	Get(id *pb.EquipId) (*pb.Equip, bool)
	GetAll() ([]*pb.Equip, error)
	GetOrCreate(equip *pb.Equip) error
	Create(equip *pb.Equip) error
	Update(equip *pb.Equip) error
	Delete(id *pb.EquipId) error
}
