package app

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.uber.org/zap"
	logger2 "p.solovev/pkg/logger"

	"github.com/go-chi/chi/v5"
	"p.solovev/internal/config"
	v1 "p.solovev/internal/controller/http/v1"
	"p.solovev/internal/usecase"
	"p.solovev/internal/usecase/cbrf"
	"p.solovev/pkg/httpserver"
)

func Run(cfg *config.Config) {

	logger2.InitLogger()

	c := usecase.NewCurrencyUseCase(cbrf.NewCurrencyReq())
	i := usecase.NewInfoUsecase(cfg.Version, cfg.Hostname)

	handler := chi.NewRouter()

	v1.NewRouter(handler, c, i)

	server := httpserver.New(handler, httpserver.Port(strconv.Itoa(cfg.Port)))
	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, os.Interrupt, syscall.SIGTERM)

	logger2.Info("Starting server", zap.Int("port", cfg.Port))

	select {
	case s := <-interruption:
		logger2.Warn("Interruption channel: ", zap.String("notification", s.String()))
	case err := <-server.Notify():
		logger2.Warn("Server notify channel: ", zap.Error(err))
	}

	err := server.Shutdown()
	if err != nil {
		logger2.Error("Error shutting down server: ", zap.Error(err))
	}
}
