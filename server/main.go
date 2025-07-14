package main

import (
	"context"
	"log"
	"net"
	"net/http"

	user "github.com/Santosh-Sohan/user-service/api"
	"github.com/Santosh-Sohan/user-service/db"
	"github.com/Santosh-Sohan/user-service/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	db.InitCassandra()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	svc := service.NewUserService()
	user.RegisterUserServiceServer(grpcServer, svc)

	go func() {
		log.Println("gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	gwMux := runtime.NewServeMux()
	err = user.RegisterUserServiceHandlerServer(context.Background(), gwMux, svc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("REST gateway on :8080")
	log.Fatal(http.ListenAndServe(":8080", gwMux))
}
