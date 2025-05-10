package app

import (
	"github.com/nextlag/in_memory_db/internal/errs"
	"github.com/nextlag/in_memory_db/internal/server/usecase/storage/engine/in_memory"
)

// createEngine initializes the proper storage engine based on configuration
func (a *App) createEngine() (*in_memory.Engine, error) {
	engineType := a.cfg.Engine.Type
	if engineType == "" {
		engineType = "in_memory"
	}

	supportedTypes := map[string]struct{}{
		"in_memory": {},
	}

	if _, ok := supportedTypes[engineType]; !ok {
		return nil, errs.ErrEngineTypeIsIncorrect
	}

	var options []in_memory.EngineOptions

	if a.cfg.Engine.PartitionsNumber > 0 {
		options = append(options, in_memory.WithPartitions(a.cfg.Engine.PartitionsNumber))
	}

	return in_memory.NewEngine(a.log, options...)
}
