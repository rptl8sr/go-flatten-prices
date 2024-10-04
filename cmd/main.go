package main

import (
	"fmt"
	"os"
	"time"

	"go-flattern-prices/app"
	"go-flattern-prices/internal/logger"
)

func main() {
	startTime := time.Now()

	a, err := app.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	a.Start()

	logger.Info(fmt.Sprintf("Готово! Время выполнения: %.2f секунд", time.Since(startTime).Seconds()))
}
