package app

import (
	"context"
	"github.com/lRhythm/gophermart/internal/client"
	"github.com/lRhythm/gophermart/internal/config"
	"github.com/lRhythm/gophermart/internal/logs"
	"github.com/lRhythm/gophermart/internal/service"
	"github.com/lRhythm/gophermart/internal/storage/pg"
	"github.com/lRhythm/gophermart/internal/transport/rest"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	logger := logs.New()
	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err)
	}
	store, err := pg.New(cfg.Database.DSN)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Close()
	accrual := client.New()
	logic := service.New(
		service.WithLogs(logger),
		service.WithConfig(&cfg.Accrual),
		service.WithStorage(store),
		service.WithClient(accrual),
	)
	defer logic.Close()

	s, err := rest.New(logger, cfg, logic)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	logic.Exchange(ctx)

	sCh := make(chan os.Signal, 1)
	signal.Notify(sCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sCh
		logger.Info("server shutting down")
		_ = s.Shutdown()
	}()

	logger.Info("server started")
	if err = s.Listen(); err != nil {
		logger.Fatal(err)
	}

	cancel()
	logger.Info("server shut down")
}
