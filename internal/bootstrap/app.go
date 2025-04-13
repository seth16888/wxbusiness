package bootstrap

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/seth16888/wxbusiness/internal/di"
	"go.uber.org/zap"
)

func StartApp(deps *di.Container) error {
	srv := InitServer(deps)
	deps.Server = srv

	// Run server
	errChan := make(chan error, 1)
	srv.Run(errChan)

	quitFunc := func() {
		srv.Shutdown()
		deps.Log.Info("Server shutdown completed")
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		quitFunc()
		return nil
	case err := <-errChan:
		deps.Log.Error("Server error", zap.Error(err))
		quitFunc()
		return err
	}
}
