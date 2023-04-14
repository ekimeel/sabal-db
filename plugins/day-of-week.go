package plugins

import (
	"encoding/json"
	"errors"
	"github.com/ekimeel/sabal-db/internal/storage"
	"github.com/ekimeel/sabal-pb/pb"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

const dayOfWeekPluginTableName = "dayOfWeek"

type DayOfWeekPlugin struct {
	data *storage.Collection[uint32, *DayOfWeek]
}

type DayOfWeek struct {
	PointId     uint32                     `json:"point-id"`
	Start       time.Time                  `json:"start"`
	End         time.Time                  `json:"end"`
	Evaluations uint32                     `json:"evaluations,omitempty"`
	Monday      storage.Statistic[float64] `json:"mon,omitempty"`
	Tuesday     storage.Statistic[float64] `json:"tue,omitempty"`
	Wednesday   storage.Statistic[float64] `json:"wed,omitempty"`
	Thursday    storage.Statistic[float64] `json:"thu,omitempty"`
	Friday      storage.Statistic[float64] `json:"fri,omitempty"`
	Saturday    storage.Statistic[float64] `json:"sat,omitempty"`
	Sunday      storage.Statistic[float64] `json:"sun,omitempty"`
}

func (d *DayOfWeekPlugin) Install() {
	d.data = storage.NewCollection[uint32, *DayOfWeek](
		dayOfWeekPluginTableName,
		0,
		storage.NewSerialKeyGenerator(0))
}

func doCalc(data []*pb.Metric, existing *storage.Statistic[float64]) {

	if len(data) == 0 {
		return
	}

	if existing == nil {
		existing = &storage.Statistic[float64]{}
	}

	var sum, sumsqrd, stddev float64
	var count int
	for i := range data {
		d := data[i]

		sum += d.Value
		if d.Value < existing.Minimum || existing.Count == 0 {
			existing.Minimum = d.Value
		}
		if d.Value > existing.Maximum || existing.Count == 0 {
			existing.Maximum = d.Value
		}
		sumsqrd += d.Value * d.Value
		count += 1
	}

	average := sum / float64(count)
	stddev = math.Sqrt(sumsqrd/float64(count) - average*average)

	if existing.Count > 0 {
		curWeight := float64(len(data)) / float64(existing.Count-len(data))
		totalWeight := curWeight - 1.0

		existing.Average = weightedAverage(existing.Average, totalWeight, average, curWeight)
		existing.StdDev = weightedAverage(existing.StdDev, totalWeight, stddev, curWeight)
		existing.Count += count

	} else {
		existing.Count = count
		existing.Average = average
		existing.StdDev = stddev
	}

}

func (d *DayOfWeekPlugin) Get(metric uint32) *DayOfWeek {
	dow, _ := d.data.Get(metric)
	return dow
}

func (d *DayOfWeekPlugin) Run(env *Environment) error {
	log.Infof("running plugin: %s", "DayOfWeek")
	metricService := env.MetricService
	pointService := env.PointService
	points, err := pointService.GetAll()

	if err != nil {
		log.Warn("failed to get points")
		return errors.New("no points")
	}

	for i := range points {
		point := points[i]

		dow, ok := d.data.Get(point.Id)

		if ok == false {
			dow = &DayOfWeek{
				PointId: point.Id,
				Start:   time.Unix(0, 0),
				End:     time.Unix(0, 0),
			}
			_, err := d.data.Create(dow)
			if err != nil {
				log.Errorf("failed to create new dayOfWeek entry: %s", err)
				return errors.New("failed to create entry")
			}
		}

		data, err := metricService.Select(point.Id, dow.End.Unix(), math.MaxInt64)
		if err != nil {
			log.Errorf("failed to read: %s", err)
			return errors.New("failed to read")
		}

		mon := make([]*pb.Metric, 0)
		tue := make([]*pb.Metric, 0)
		wed := make([]*pb.Metric, 0)
		thr := make([]*pb.Metric, 0)
		fri := make([]*pb.Metric, 0)
		sat := make([]*pb.Metric, 0)
		sun := make([]*pb.Metric, 0)

		for j := range data {
			d := data[j]
			switch d.Timestamp.AsTime().Weekday() {
			case time.Monday:
				mon = append(mon, d)
			case time.Tuesday:
				tue = append(tue, d)
			case time.Wednesday:
				wed = append(wed, d)
			case time.Thursday:
				thr = append(thr, d)
			case time.Friday:
				fri = append(fri, d)
			case time.Saturday:
				sat = append(sat, d)
			case time.Sunday:
				sun = append(sun, d)
			}
		}

		doCalc(mon, &dow.Monday)
		doCalc(tue, &dow.Tuesday)
		doCalc(wed, &dow.Wednesday)
		doCalc(thr, &dow.Thursday)
		doCalc(fri, &dow.Friday)
		doCalc(sat, &dow.Saturday)
		doCalc(sun, &dow.Sunday)

	}

	return nil

}

func (d *DayOfWeek) String() string {
	s, err := json.Marshal(d)
	if err != nil {
		return "error"
	}
	return string(s)
}
