package app

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"

	"github.com/nextlag/in_memory_db/config"
	"github.com/nextlag/in_memory_db/internal/errs"
	"github.com/nextlag/in_memory_db/internal/server/launcher"
	"github.com/nextlag/in_memory_db/internal/server/repository"
	"github.com/nextlag/in_memory_db/internal/server/usecase"
	"github.com/nextlag/in_memory_db/internal/server/usecase/compute"
	"github.com/nextlag/in_memory_db/internal/server/usecase/storage/engine/in_memory"
)

const fileConfig = "config.json"

// App represents the main application container
type App struct {
	eng *in_memory.Engine // Storage engine
	srv *launcher.Runtime // Server runtime
	wg  sync.WaitGroup    // WaitGroup for graceful shutdown
	cfg *config.Config    // Application configuration
	log *slog.Logger      // Logger instance
}

// New creates a new application instance
func New() (*App, error) {
	cfg, err := config.New(fileConfig)
	if err != nil {
		return nil, errors.Join(errs.ErrFailedInitConfig, err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	srv, err := launcher.New(cfg, log)
	if err != nil {
		return nil, errors.Join(errs.ErrFailedInitLauncher, err)
	}

	a := &App{
		cfg: cfg,
		log: log,
		srv: srv,
	}

	a.eng, err = a.createEngine()
	if err != nil {
		return nil, errors.Join(errs.ErrFailedInitEngine, err)
	}

	return a, nil
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	comp, err := compute.New(a.log)
	if err != nil {
		return errors.Join(errs.ErrFailedInitCompute, err)
	}

	repo, err := repository.New(a.eng, a.log)
	if err != nil {
		return errors.Join(errs.ErrFailedInitStorage, err)
	}

	uc, err := usecase.New(comp, repo, a.log)
	if err != nil {
		return errors.Join(errs.ErrFailedInitUseCase, err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	a.wg.Add(1)

	go func() {
		defer a.wg.Done()

		if err = a.srv.Runtime(ctx, func(ctx context.Context, query []byte) []byte {
			response := uc.HandleQuery(ctx, string(query))

			return []byte(response)
		}); err != nil {
			a.log.Error("server runtime failed", errs.ErrLog(err))

			cancel()
		}
	}()

	<-ctx.Done()

	a.srv.Close()
	a.wg.Wait()

	a.log.Info("Close completed")

	return nil
}
