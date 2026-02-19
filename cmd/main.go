package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/delyke/tasks_and_commands_service/closer"
	"github.com/delyke/tasks_and_commands_service/internal/app"
	"github.com/delyke/tasks_and_commands_service/internal/config"
	"github.com/delyke/tasks_and_commands_service/logger"
)

const configPath = "./deploy/env/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		log.Println(err)
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		log.Println(err)
		return
	}
}

// gracefulShutdown выполняет корректное завершение работы сервиса.
//
// Создаёт новый контекст с таймаутом 5 секунд для закрытия всех зависимостей.
func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
