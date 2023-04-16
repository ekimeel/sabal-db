package main

import (
	"context"
	"github.com/ekimeel/sabal-pb/pb"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat/distuv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	equipClient := pb.NewEquipServiceClient(conn)
	pointClient := pb.NewPointServiceClient(conn)
	metricClient := pb.NewMetricServiceClient(conn)

	equip := &pb.Equip{
		Name: "my-equip",
	}

	equip, err = equipClient.Create(context.Background(), equip)
	if err != nil {
		panic(err)
	}

	_, _ = pointClient.Create(context.Background(), &pb.Point{Name: "point-a", EquipId: equip.Id})
	_, _ = pointClient.Create(context.Background(), &pb.Point{Name: "point-b", EquipId: equip.Id})
	_, _ = pointClient.Create(context.Background(), &pb.Point{Name: "point-c", EquipId: equip.Id})
	_, _ = pointClient.Create(context.Background(), &pb.Point{Name: "point-d", EquipId: equip.Id})

	equipList, err := equipClient.GetAll(context.Background(), &pb.ListRequest{Limit: 1000, Offset: 0})
	if err != nil {
		panic(err)
	}

	log.Printf("equipList: %v", equipList)

	pointList, err := pointClient.GetAllByEquip(context.Background(), &pb.EquipId{Id: equip.Id})
	if err != nil {
		panic(err)
	}

	log.Printf("pointList: %v", pointList)

	clock := time.Now()

	writeBuffer := &pb.MetricList{}
	writeBuffer.Metrics = make([]*pb.Metric, 0, 5000)

	for i := 0; i < (40000); i++ {
		clock = clock.Add(1 * time.Minute)

		for _, point := range pointList.Points {
			value := distuv.Normal{Mu: 55, Sigma: 5}

			writeBuffer.Metrics = append(writeBuffer.Metrics, &pb.Metric{
				PointId:   point.Id,
				Value:     value.Rand(),
				Timestamp: &timestamp.Timestamp{Seconds: clock.Unix()}})

		}

		if len(writeBuffer.Metrics) >= 5000 {
			start := time.Now()
			ok, err := metricClient.WriteList(context.Background(), writeBuffer)
			log.Printf("write: %v, %s", len(writeBuffer.Metrics), time.Since(start))
			writeBuffer.Metrics = nil
			writeBuffer.Metrics = make([]*pb.Metric, 0, 5000)

			if err != nil {
				panic(err)
			}
			if ok.Accepted == false {
				panic("failed to accept write")
			}
		}

	}
	start := time.Now()
	ok, err := metricClient.WriteList(context.Background(), writeBuffer)
	log.Printf("write: %v, %s", len(writeBuffer.Metrics), time.Since(start))
	if err != nil {
		panic(err)
	}
	if ok.Accepted == false {
		panic("failed to accept write")
	}

}
