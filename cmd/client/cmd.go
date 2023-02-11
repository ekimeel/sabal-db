package main

import (
	"context"
	"github.com/ekimeel/sabal-db/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	equipClient := pb.NewEquipServiceClient(conn)
	pointClient := pb.NewPointServiceClient(conn)

	equip := &pb.Equip{
		Uuid:    "t.abc",
		Enabled: true,
	}

	equipClient.Create(context.Background(), equip)
	pointClient.Create(context.Background(), &pb.Point{Name: "rpm", EquipUuid: "t.abc", Enabled: true})
	pointClient.Create(context.Background(), &pb.Point{Name: "speed", EquipUuid: "t.abc", Enabled: true})
	pointClient.Create(context.Background(), &pb.Point{Name: "temp", EquipUuid: "t.abc", Enabled: true})
	pointClient.Create(context.Background(), &pb.Point{Name: "oilPressure", EquipUuid: "t.xxx", Enabled: true})

	equipList, err := equipClient.GetAll(context.Background(), &pb.ListRequest{Limit: 1000, Offset: 0})
	if err != nil {
		panic(err)
	}

	log.Printf("equipList: %v", equipList)

	pointList, err := pointClient.GetAllByEquipUUID(context.Background(), &pb.EquipUUID{Id: "t.abc"})
	if err != nil {
		panic(err)
	}

	log.Printf("pointList: %v", pointList)
}
