package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq" // PostgreSQL driver
	"google.golang.org/grpc"

	advert "github.com/s21platform/advert-proto/advert-proto"
	"github.com/s21platform/advert-service/internal/config"
	"github.com/s21platform/advert-service/internal/infra"
	db "github.com/s21platform/advert-service/internal/repository/postgres"
	"github.com/s21platform/advert-service/internal/service"
)

func main() {
	cfg := config.MustLoad()

	dbRepo := db.New(cfg)
	defer dbRepo.Close()

	advertService := service.New(dbRepo)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
		),
	)

	advert.RegisterAdvertServiceServer(server, advertService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("cannot listen port: %s; Error: %v", cfg.Service.Port, err)
	}
	if err = server.Serve(lis); err != nil {
		log.Fatalf("cannot start grpc, port: %s; Error: %v", cfg.Service.Port, err)
	}
}
