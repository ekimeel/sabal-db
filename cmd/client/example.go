package main

import (
	"context"
	"github.com/ekimeel/sabal-pb/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	pointA, _ := pointClient.Create(context.Background(), &pb.Point{Name: "point-a", EquipId: equip.Id})
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

	for i := 0; i < 100000; i++ {
		r, err := metricClient.Write(
			context.Background(),
			&pb.Metric{PointId: pointA.Id, Timestamp: timestamppb.Now(), Value: float64(i)})
		if err != nil {
			panic(err)
		}

		last, err := metricClient.Poll(context.Background(), &pb.PointId{Id: pointA.Id})
		if last != nil {
			log.Printf("r: %v, poll:%f, ts:%s", r, last.Value, last.Timestamp)
		}

		time.Sleep(10 * time.Millisecond)
	}

}
