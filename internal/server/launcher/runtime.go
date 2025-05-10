package launcher

import (
	"context"
	"log/slog"

	"github.com/nextlag/in_memory_db/config"
	"github.com/nextlag/in_memory_db/internal/errs"
	"github.com/nextlag/in_memory_db/internal/server/entity"
	"github.com/nextlag/in_memory_db/internal/server/launcher/cmd"
)

const (
	cmdServer = "cmd"
)

type launcher interface {
	ServerRun(ctx context.Context, handler entity.Handler) error
	Close()
}

type Runtime struct {
	srv launcher
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) (*Runtime, error) {
	if cfg == nil {
		return nil, errs.ErrConfigIsNil
	}

	if log == nil {
		return nil, errs.ErrLoggerIsNil
	}

	var srv launcher

	switch cfg.Server.Type {
	case cmdServer:
		srv = cmd.New(log)
	default:
		return nil, errs.ErrUnknownServerType
	}

	return &Runtime{
		srv: srv,
		log: log,
	}, nil
}

func (r *Runtime) Runtime(ctx context.Context, handler entity.Handler) error {
	return r.srv.ServerRun(ctx, handler)
}

func (r *Runtime) Close() {
	r.srv.Close()
}
