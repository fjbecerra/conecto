package main

import (
	"conecto/factories"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	appConfig, err := factories.LoadConfig[factories.AppConfig](
		"./config/conecto.json",
	)

	if err != nil {
		panic(err)
	}


	runner := factories.NewConecto(appConfig.ConectoConfig).Build()


	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	slog.Info("application starting")


	go runner.RunWorker(ctx)

	go runner.RunScheduler(ctx)


	server := &http.Server{
		Addr: fmt.Sprintf(
			":%s",
			strconv.Itoa(appConfig.ConectoConfig.HttpServerConfig.Port),
		),
		Handler: runner.RunRouter(),
	}


	go func() {

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatalf(
				"server error: %v",
				err,
			)
		}
	}()


	<-ctx.Done()


	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancelShutdown()

	// 1. Stop accepting new HTTP requests and wait for existing ones
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error(
			"server shutdown failed",
			"error",
			err,
		)
	}

	// 2. Close database connections and other resources
	if err := runner.Closed(); err != nil {
		slog.Error(
			"resource cleanup failed",
			"error",
			err,
		)
	}
}

