package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/willtowle1/parkn/internal/common/logger"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type App struct {
	logger logger.Logger
	Server Server
}

func NewApp(logger logger.Logger, server *http.Server) *App {
	return &App{
		logger: logger,
		Server: server,
	}
}

func (a *App) Start(ctx context.Context, errs chan error, serverAddress string) {
	go func() {
		errs <- a.Server.ListenAndServe()
	}()

	a.logger.Info(ctx, fmt.Sprintf("parkn-service running on address: %s", serverAddress))
}

func (a *App) Shutdown(ctx context.Context, gracePeriod time.Duration) error {
	termCtx, cancel := context.WithTimeout(ctx, gracePeriod)
	defer cancel()
	err := a.Server.Shutdown(termCtx)
	return err
}

func WaitForTermination(ctx context.Context, logger logger.Logger, errs chan error) {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errs:
		logger.Error(ctx, "shutting down caused by error", err)
	case sig := <-signals:
		logger.Error(ctx, "shutting down from signal", fmt.Errorf("signal: %s", sig.String()))
	}
}
