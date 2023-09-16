package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SergeiSkv/pow/server/assets"
	"github.com/SergeiSkv/pow/server/config"
	"github.com/SergeiSkv/pow/server/handlers"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		logrus.Fatalf("parse config err %v", err)
	}
	fmt.Printf("Server running on %s\n", cfg.Host())

	quotes, err := assets.GetQuotes()
	if err != nil {
		logrus.Fatalf("read quotes %v", err)
	}

	h := handlers.NewHandler(cfg.TargetPrefix, cfg.ChallengeLength, quotes)

	stop := setupSignalHandling()

	go runServer(context.Background(), h, cfg.Host(), stop)

	<-stop
	fmt.Println("Received signal, shutting down...")
}

func setupSignalHandling() chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	return stop
}

func runServer(ctx context.Context, h handlers.Handler, host string, stop chan os.Signal) {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		logrus.Fatalf("Error listening: %s", err.Error())
	}
	defer listen.Close()

	for {
		select {
		case <-stop:
			return
		default:
			conn, err := listen.Accept()
			if err != nil {
				logrus.Errorf("Error accepting connection: %s", err.Error())
				continue
			}

			ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
			if err = h.Handle(ctx, conn); err != nil {
				logrus.Error(err)
			}
			cancel()
		}
	}
}
