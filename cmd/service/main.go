package main

import (
	"fmt"

	"github.com/s21platform/advert-service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	// logger := logger_lib.New(cfg.Service.Name, cfg.Platform.Env)
	fmt.Print(cfg)
}
