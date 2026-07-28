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


	go runner.Worker.Run(ctx)

	go runner.Scheduler.Run(ctx)


	server := &http.Server{
		Addr: fmt.Sprintf(
			":%s",
			strconv.Itoa(appConfig.ConectoConfig.HttpServerConfig.Port),
		),
		Handler: runner.Router,
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


	server.Shutdown(shutdownCtx)
}

