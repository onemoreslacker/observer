package scrapperservice

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/application/scrapper/core"
)

type ScrapperService struct {
	scr *core.Scrapper
	srv *http.Server
}

func New(scr *core.Scrapper, srv *http.Server) (*ScrapperService, error) {
	return &ScrapperService{
		scr: scr,
		srv: srv,
	}, nil
}

func (s *ScrapperService) Run() error {
	srvErr := make(chan error, 1)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
	}()

	scrapperErr := make(chan error, 1)

	go func() {
		if err := s.scr.Run(); err != nil {
			scrapperErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErr:
		slog.Error(
			"server error",
			slog.String("msg", err.Error()),
			slog.String("service", "scrapper"),
		)
	case err := <-scrapperErr:
		slog.Error(
			"scrapper error",
			slog.String("msg", err.Error()),
			slog.String("service", "scrapper"),
		)
	case sig := <-stop:
		slog.Info(
			"received shutdown signal",
			slog.String("signal", sig.String()),
			slog.String("service", "scrapper"),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
