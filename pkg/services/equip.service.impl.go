package services

import (
	"github.com/ekimeel/sabal-db/internal/storage"
	"github.com/ekimeel/sabal-pb/pb"
	log "github.com/sirupsen/logrus"
	"sync"
)

var equipServiceInstance *EquipServiceImpl
var equipServiceOnce sync.Once

type EquipServiceImpl struct {
	data *storage.Collection[uint32, *pb.Equip]
}

func GetEquipService() EquipService {

	equipServiceOnce.Do(func() {
		log.Tracef("creating [%s] once", "equipServiceInstance")

		equipServiceInstance = &EquipServiceImpl{}
		equipServiceInstance.data = storage.NewCollection[uint32, *pb.Equip]("equip", 0,
			storage.NewSerialKeyGenerator(0))
	})

	return equipServiceInstance
}

func (s *EquipServiceImpl) Create(equip *pb.Equip) error {
	log.Tracef("Create: %s", equip.String())

	id, err := s.data.Create(equip)
	if err != nil {
		log.Errorf("failed to create point: %s", err)
		return err
	}

	equip.Id = id
	return nil
}

func (s *EquipServiceImpl) Update(equip *pb.Equip) error {
	log.Tracef("Update: %s", equip.String())

	err := s.data.Replace(equip.Id, equip)
	if err != nil {
		return err
	}

	return nil
}

func (s *EquipServiceImpl) Get(uuid *pb.EquipId) (*pb.Equip, bool) {
	log.Infof("Get: %d", uuid.GetId())
	return s.data.Get(uuid.GetId())
}

func (s *EquipServiceImpl) GetOrCreate(equip *pb.Equip) error {
	log.Infof("GetOrCreate: %s", equip.String())

	existing, loaded := s.Get(&pb.EquipId{Id: equip.GetId()})

	if loaded {
		equip = existing
	} else {
		err := s.Create(equip)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *EquipServiceImpl) GetAll() ([]*pb.Equip, error) {
	log.Info("GetAll")
	return s.data.Values(), nil
}

func (s *EquipServiceImpl) Delete(id *pb.EquipId) error {
	log.Infof("delete equip: %d", id.GetId())
	return s.data.Remove(id.GetId())
}
