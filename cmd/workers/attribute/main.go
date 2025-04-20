package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"

	kafka_lib "github.com/s21platform/kafka-lib"
	pkg "github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/advert-service/internal/config"
	new_attribute "github.com/s21platform/advert-service/internal/databus/new_attribute"
	db "github.com/s21platform/advert-service/internal/repository/postgres"
)

func main() {
	cfg := config.MustLoad()

	dbRepo := db.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Println("failed to connect graphite: ", err)
	}

	kafkaConfig := kafka_lib.DefaultConsumerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.SetAttributeTopic, cfg.Kafka.Group)

	consumer, err := kafka_lib.NewConsumer(kafkaConfig, metrics)
	if err != nil {
		log.Fatal("failed to create kafka consumer: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.RegisterHandler(ctx, new_attribute.NewAttribute)

	fmt.Println("Consumer started")

	<-ctx.Done()
}
