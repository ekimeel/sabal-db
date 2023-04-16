package time_quality

import (
	"encoding/json"
	"errors"
	"github.com/ekimeel/sabal-db/internal/storage"
	"github.com/ekimeel/sabal-db/plugins"
	"github.com/ekimeel/sabal-pb/pb"
	log "github.com/sirupsen/logrus"
	"math"
	"sort"
	"time"
)

const qualityPluginTableName = "time-quality"

var timeQualityPluginImpl *TimeQualityPlugin

type TimeQualityPlugin struct {
	data *storage.Collection[uint32, *TimeQuality]
}

type TimeQuality struct {
	PointId                            uint32    `json:"point-id"`
	Start                              time.Time `json:"start"`
	End                                time.Time `json:"end"`
	Count                              int64     `json:"count"`
	MeanTimeBetweenObservations        float64   `json:"mean-time-between-observations"`
	MaxTimeBetweenObservations         int64     `json:"max-time-between-observations"`
	MinTimeBetweenObservations         int64     `json:"min-time-between-observations"`
	TimeIntervalCoefficientOfVariation float64   `json:"time-interval-coefficient-of-variation"`
	TimeIntervalStandardDeviation      float64   `json:"time-interval-standard-deviation"`
	FillFactor                         uint32    `json:"fill-factor"`
	Score                              uint32    `json:"score"`
}

func Install() {
	timeQualityPluginImpl = &TimeQualityPlugin{}
	timeQualityPluginImpl.data = storage.NewCollection[uint32, *TimeQuality](
		qualityPluginTableName,
		0,
		storage.NewSerialKeyGenerator(0))
}

func computeScore(existing *TimeQuality) {
	// compute fill factor
	span := float64(existing.End.Unix() - existing.Start.Unix())
	mean := existing.MeanTimeBetweenObservations
	count := float64(existing.Count)

	fillFactor := count / (span / mean) * 100
	cv := 100 - existing.TimeIntervalCoefficientOfVariation

	score := (fillFactor + cv) / 2

	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	existing.Score = uint32(score)

}

func doQualityCalc(data []*pb.Metric, existing *TimeQuality) {

	if len(data) < 4 || existing == nil {
		return
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp.Seconds < data[j].Timestamp.Seconds
	})

	intervals := make([]int64, len(data)-1)

	// calculate time intervals and mean
	var sum int64
	for i := 1; i < len(data); i++ {
		intervals[i-1] = data[i].Timestamp.Seconds - data[i-1].Timestamp.Seconds

		sum += intervals[i-1]
		if intervals[i-1] > existing.MaxTimeBetweenObservations {
			existing.MaxTimeBetweenObservations = intervals[i-1]
		}
		if intervals[i-1] < existing.MinTimeBetweenObservations {
			existing.MinTimeBetweenObservations = intervals[i-1]
		}
	}
	mean := float64(sum) / float64(len(intervals))

	// compute standard deviation and coefficient of variation
	var sd, cv float64
	for _, interval := range intervals {
		diff := float64(interval) - mean
		sd += diff * diff
	}
	sd = math.Sqrt(sd / float64(len(intervals)-1))
	cv = sd / mean * 100.0

	// compute fill factor
	span := data[len(data)-1].Timestamp.Seconds - data[0].Timestamp.Seconds
	fillFactor := float64(len(data)) / (float64(span) / mean) * 100

	if existing.Count > 0 {
		//weighted average
		curWeight := float64(len(data)) / float64(existing.Count)
		totalWeight := curWeight - 1.0

		existing.MeanTimeBetweenObservations = weightedAverage(
			existing.MeanTimeBetweenObservations, totalWeight,
			mean, curWeight)

		existing.TimeIntervalStandardDeviation = weightedAverage(
			existing.TimeIntervalStandardDeviation, totalWeight,
			sd, curWeight)

		existing.TimeIntervalCoefficientOfVariation = weightedAverage(
			existing.TimeIntervalCoefficientOfVariation, totalWeight,
			cv, curWeight)

		existing.FillFactor = uint32(weightedAverage(
			float64(existing.FillFactor), totalWeight,
			fillFactor, curWeight))

		if existing.Start.Unix() > data[0].Timestamp.Seconds {
			existing.Start = data[0].Timestamp.AsTime()
		}

		if existing.End.Unix() < data[len(data)-1].Timestamp.Seconds {
			existing.End = data[len(data)-1].Timestamp.AsTime()
		}

	} else {
		existing.MeanTimeBetweenObservations = mean
		existing.TimeIntervalStandardDeviation = sd
		existing.TimeIntervalCoefficientOfVariation = cv
		existing.Start = data[0].Timestamp.AsTime()
		existing.End = data[len(data)-1].Timestamp.AsTime()
		existing.FillFactor = uint32(fillFactor)
	}

	existing.Count += int64(len(data))
	computeScore(existing)

}

func (d *TimeQualityPlugin) Get(metric uint32) *TimeQuality {
	dow, _ := d.data.Get(metric)
	return dow
}

func Run(env *plugins.Environment) error {
	d := timeQualityPluginImpl

	log.Infof("running plugin: %s", "TimeQualityPlugin")
	metricService := env.MetricService
	pointService := env.PointService

	points, err := pointService.GetAll()

	if err != nil {
		log.Warn("failed to get points")
		return errors.New("no points")
	}

	for i := range points {
		point := points[i]

		cur, ok := d.data.Get(point.Id)

		if ok == false {
			cur = &TimeQuality{
				PointId:                    point.Id,
				Start:                      time.Unix(0, 0),
				End:                        time.Unix(0, 0),
				MaxTimeBetweenObservations: math.MinInt64,
				MinTimeBetweenObservations: math.MaxInt64,
			}
			_, err := d.data.Create(cur)
			if err != nil {
				log.Errorf("failed to create new dayOfWeek entry: %s", err)
				return errors.New("failed to create entry")
			}
		}

		data, err := metricService.Select(point.Id, cur.End.Unix(), math.MaxInt64)
		if err != nil {
			log.Errorf("failed to read: %s", err)
			return errors.New("failed to read")
		}

		doQualityCalc(data, cur)

	}

	return nil

}

func weightedAverage(value1, weight1, value2, weight2 float64) float64 {
	return (value1*weight1 + value2*weight2) / (weight1 + weight2)
}

func (d *TimeQuality) String() string {
	s, err := json.Marshal(d)
	if err != nil {
		return "error"
	}
	return string(s)
}
