package bootstrap

import (
	"fmt"
	"github.com/ekimeel/db-api/internal/env"
	"github.com/ekimeel/db-api/pb"
	pbapi2 "github.com/ekimeel/db-api/pkg/pbapi"
	services2 "github.com/ekimeel/db-api/pkg/services"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"net"
	"os"
)

func Start() {
	fmt.Printf("bootstraping... \n")

	step := 0
	setupLogging()

	log.Info("******************************************************")
	log.Info("Application - device-simulator")
	log.Info("")
	log.Info("(c) Michael J. Lee")
	log.Info("Warning: Unauthorized use of this system is prohibited")
	log.Info("*******************************************************")

	step++
	log.Infof("%d: loading environment \n", step)
	setupEnv()

	step++
	log.Infof("%d: loading rpc server \n", step)
	setupGrpc()

}

func setupLogging() {

	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	log.SetOutput(os.Stdout)
	logLevel, err := log.ParseLevel(env.GetLogLevel())
	if err != nil {
		log.Errorf("failed to parse log level: %s", err)
		log.Warnf("no %s set, defaulting to %s", logLevel, log.TraceLevel.String())
		logLevel = log.TraceLevel

	}

	log.SetLevel(logLevel)
}

func setupEnv() {

	log.Tracef("loading config file [%s]", env.GetConfigFile())

	var config env.Config
	configFile, err := os.ReadFile(env.GetConfigFile())
	if err != nil {
		log.Fatalf("failed to load config file: %s", err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("failed to unmarshal config file: %s", err)
	}

	env.SetConfig(&config)

}

func setupGrpc() {

	equipService := services2.NewEquipService()
	equipServer, err := pbapi2.NewGrpcEquipServer(equipService)
	if err != nil {
		log.Fatal("cannot create grpc postServer: ", err)
	}

	pointService := services2.NewPointService()
	pointServer, err := pbapi2.NewGrpcPointServer(pointService)
	if err != nil {
		log.Fatal("cannot create grpc postServer: ", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterEquipServiceServer(grpcServer, equipServer)
	pb.RegisterPointServiceServer(grpcServer, pointServer)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

}
