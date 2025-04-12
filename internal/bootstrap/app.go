package bootstrap

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/seth16888/wxbusiness/internal/di"
	"go.uber.org/zap"
)

func StartApp() error {
	srv := InitServer(di.Get().Conf, di.Get().Log)
	di.Get().Server = srv

	// Run server
	errChan := make(chan error, 1)
	srv.Run(errChan)

	quitFunc := func() {
		srv.Shutdown()
		di.Get().Log.Info("Server shutdown completed")
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		quitFunc()
		return nil
	case err := <-errChan:
		di.Get().Log.Error("Server error", zap.Error(err))
		quitFunc()
		return err
	}
}
