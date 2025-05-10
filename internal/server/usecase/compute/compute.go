package compute

import (
	"errors"
	"log/slog"

	"github.com/nextlag/in_memory_db/internal/errs"
)

type Compute struct {
	log *slog.Logger
}

func New(log *slog.Logger) (*Compute, error) {
	if log == nil {
		return nil, errs.ErrLoggerIsNil
	}
	return &Compute{
		log: log,
	}, nil
}

func (c *Compute) Pipeline(queryStr string) (Query, error) {
	tokens, err := c.Parse(queryStr)
	if err != nil {
		return Query{}, errors.Join(errs.ErrParse, err)
	}

	query, err := c.Analyze(tokens)
	if err != nil {
		return Query{}, errors.Join(errs.ErrAnalyze, err)
	}

	return query, nil
}
