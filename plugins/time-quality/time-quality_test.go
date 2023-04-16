package time_quality

import (
	"fmt"
	"github.com/ekimeel/sabal-db/pkg/services"
	"github.com/ekimeel/sabal-db/plugins"
	"github.com/ekimeel/sabal-pb/pb"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/stat/distuv"
	"testing"
	"time"
)

func TestQualityPlugin_Run(t *testing.T) {

	plugin := TimeQualityPlugin{}
	plugin.Install()

	equipService := services.GetEquipService()
	pointService := services.GetPointService()
	metricService := services.GetMetricService()

	equip := &pb.Equip{Name: "equip-a"}
	err := equipService.GetOrCreate(equip)
	assert.Nil(t, err)

	pointA := &pb.Point{Name: "point-a", EquipId: equip.Id}
	pointB := &pb.Point{Name: "point-b", EquipId: equip.Id}
	pointC := &pb.Point{Name: "point-c", EquipId: equip.Id}
	pointD := &pb.Point{Name: "point-d", EquipId: equip.Id}
	pointE := &pb.Point{Name: "point-e", EquipId: equip.Id}

	_ = pointService.Create(pointA)
	_ = pointService.Create(pointB)
	_ = pointService.Create(pointC)
	_ = pointService.Create(pointD)
	_ = pointService.Create(pointE)

	points, err := pointService.GetAll()
	assert.Nil(t, err)

	count := 0
	clock := time.Now()

	for i := 0; i < (40000); i++ {
		clock = clock.Add(1 * time.Minute)
		for _, point := range points {
			value := distuv.Normal{Mu: 55, Sigma: 5}
			ok, err := metricService.Write(&pb.Metric{
				PointId:   point.Id,
				Value:     value.Rand(),
				Timestamp: &timestamp.Timestamp{Seconds: clock.Unix()}})
			assert.True(t, ok)
			assert.Nil(t, err)
			count += 1
		}
	}
	fmt.Printf("count: %d", count)
	_ = metricService.Flush()
	err = plugin.Run(&plugins.Environment{
		MetricService: services.GetMetricService(),
		PointService:  services.GetPointService(),
		EquipService:  services.GetEquipService()})

	for _, point := range points {
		q := plugin.Get(point.Id)
		s := fmt.Sprintf(
			"point:%s, \n\tdata:%s \n",
			point.Name,
			q.String())
		fmt.Println(s)
	}
}
