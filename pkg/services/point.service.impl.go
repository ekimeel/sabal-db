package services

import (
	"github.com/ekimeel/sabal-db/internal/storage"
	"github.com/ekimeel/sabal-pb/pb"
	log "github.com/sirupsen/logrus"
	"sync"
)

type PointServiceImpl struct {
	data *storage.Collection[uint32, *pb.Point]
}

var pointServiceInstance *PointServiceImpl
var pointServiceOnce sync.Once

func GetPointService() PointService {

	pointServiceOnce.Do(func() {
		log.Tracef("creating [%s] once", "pointServiceInstance")

		pointServiceInstance = &PointServiceImpl{}
		pointServiceInstance.data = storage.NewCollection[uint32, *pb.Point]("point", 0,
			storage.NewSerialKeyGenerator(0))
	})

	return pointServiceInstance
}

func (s *PointServiceImpl) Create(point *pb.Point) error {
	log.Tracef("Create: %s", point.String())

	id, err := s.data.Create(point)
	if err != nil {
		log.Errorf("failed to create point: %s", err)
		return err
	}

	point.Id = id
	return nil
}

func (s *PointServiceImpl) Update(point *pb.Point) error {
	log.Tracef("Update: %s", point.String())

	err := s.data.Replace(point.Id, point)
	if err != nil {
		return err
	}

	return nil
}

func (s *PointServiceImpl) Get(id *pb.PointId) (*pb.Point, bool) {
	log.Tracef("Get: %s", id.String())
	return s.data.Get(id.GetId())
}

func (s *PointServiceImpl) GetAll() ([]*pb.Point, error) {
	log.Tracef("GetAll")
	return s.data.Values(), nil
}

func (s *PointServiceImpl) GetOrCreate(point *pb.Point) error {
	log.Infof("GetOrCreate: %s", point.String())

	existing, loaded := s.Get(&pb.PointId{Id: point.GetId()})

	if loaded {
		point = existing
	} else {
		err := s.Create(point)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PointServiceImpl) Delete(id *pb.PointId) error {
	log.Tracef("Delete: %s", id.String())
	return s.data.Remove(id.GetId())
}

func (s *PointServiceImpl) GetAllByEquip(id *pb.EquipId) ([]*pb.Point, error) {
	log.Tracef("GetAllByEquip: %s", id.String())
	result := s.data.Filter(func(s uint32, point *pb.Point) bool {
		return point.EquipId == id.Id
	})

	return result, nil
}
