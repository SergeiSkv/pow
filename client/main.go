package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SergeiSkv/pow/client/config"
	"github.com/SergeiSkv/pow/client/executor"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		logrus.Fatalf("parse config err %v", err)
	}

	run(cfg)
}

func run(cfg *config.Config) {
	shutdownCtx, cancel := setupSignalHandling()
	defer cancel() // Ensure all paths release resources

	mainLoop(shutdownCtx, cfg)
}

func setupSignalHandling() (context.Context, context.CancelFunc) {
	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Shutdown context to stop all ongoing operations when interrupted
	shutdownCtx, cancel := context.WithCancel(context.Background())

	go func() {
		<-stop
		logrus.Info("Received signal, shutting down...")
		cancel()
	}()

	return shutdownCtx, cancel
}

func mainLoop(ctx context.Context, cfg *config.Config) {
	for {
		select {
		case <-ctx.Done():
			logrus.Info("Shutdown complete.")
			return
		default:
			ex, err := executor.NewExecutor(cfg.ServerAddress(), cfg.TargetPrefix)
			if err != nil {
				logrus.Fatalf("init executor %v", err)
			}
			res, err := ex.Execute(context.Background())
			if err != nil {
				logrus.Error(err)
				time.Sleep(1 * time.Second) // Optional: rate limiting
				continue
			}

			logrus.Info(res)
			time.Sleep(1 * time.Second) // Optional: rate limiting
		}
	}
}
