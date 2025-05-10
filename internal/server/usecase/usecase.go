package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/nextlag/in_memory_db/internal/errs"
	"github.com/nextlag/in_memory_db/internal/server/entity"
	"github.com/nextlag/in_memory_db/internal/server/repository"
	"github.com/nextlag/in_memory_db/internal/server/usecase/compute"
)

// ComputeUseCase defines the interface for query processing
type ComputeUseCase interface {
	Pipeline(queryStr string) (compute.Query, error)
}

// UseCase orchestrates the business logic between compute and data layers
type UseCase struct {
	compute ComputeUseCase
	repo    repository.DataStore
	log     *slog.Logger
}

// New creates a new UseCase instance with required dependencies
func New(compute ComputeUseCase, repo repository.DataStore, log *slog.Logger) (*UseCase, error) {
	if compute == nil {
		return nil, errs.ErrComputeIsNil
	}

	if repo == nil {
		return nil, errs.ErrStorageIsNil
	}

	if log == nil {
		return nil, errs.ErrLoggerIsNil
	}

	return &UseCase{
		compute: compute,
		repo:    repo,
		log:     log,
	}, nil
}

func (uc *UseCase) HandleQuery(ctx context.Context, queryStr string) string {
	uc.log.Debug("handling query", slog.String("query", queryStr))

	query, err := uc.compute.Pipeline(queryStr)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	switch query.CommandID() {
	case compute.SetCommandID:
		return uc.handleSetQuery(ctx, query)
	case compute.GetCommandID:
		return uc.handleGetQuery(ctx, query)
	case compute.DelCommandID:
		return uc.handleDelQuery(ctx, query)
	default:
		uc.log.Error("compute layer is incorrect", slog.Int("command_id", query.CommandID()))
	}

	return entity.ResponseErr + "internal error"
}

func (uc *UseCase) handleSetQuery(ctx context.Context, query compute.Query) string {
	args := query.Arguments()

	if err := uc.repo.Set(ctx, args[0], args[1]); err != nil {
		return fmt.Sprintf("%s %s", entity.ResponseErr, err.Error())
	}

	return entity.ResponseOk
}

func (uc *UseCase) handleGetQuery(ctx context.Context, query compute.Query) string {
	args := query.Arguments()

	value, err := uc.repo.Get(ctx, args[0])
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return entity.ResponseNotFound
		}

		return fmt.Sprintf("%s %s", entity.ResponseErr, err.Error())
	}

	return value
}

func (uc *UseCase) handleDelQuery(ctx context.Context, query compute.Query) string {
	args := query.Arguments()

	err := uc.repo.Del(ctx, args[0])
	if errors.Is(err, errs.ErrNotFound) {
		return entity.ResponseNotFound
	}

	if err != nil {
		return fmt.Sprintf("%s %s", entity.ResponseErr, err.Error())
	}

	return entity.ResponseOk
}
