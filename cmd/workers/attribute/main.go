package attribute

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kafka_lib "github.com/s21platform/kafka-lib"
	pkg "github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/advert-service/internal/config"
	"github.com/s21platform/advert-service/internal/repository/postgres"
)

func main() {
	cfg := config.MustLoad()

	dbRepo := postgres.New(cfg)
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

	consumer.RegisterHandler(ctx, func(ctx context.Context, msg []byte) error {
		// Пример десериализации JSON сообщения
		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}

		// Обработка сообщения
		fmt.Printf("Received message: %v\n", data)
		return nil
	})

	fmt.Println("Consumer started")

	time.Sleep(time.Second)
}
