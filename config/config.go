package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/nextlag/in_memory_db/internal/errs"
)

type (
	Config struct {
		Server  *Server  `json:"launcher"`
		Engine  *Engine  `json:"engine"`
		Logging *Logging `json:"logging"`
	}

	Server struct {
		Type string `json:"type"`
	}

	Engine struct {
		Type             string `yaml:"type"`
		PartitionsNumber uint   `yaml:"partitions_number"`
	}

	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	}
)

func New(fileConfig string) (*Config, error) {
	if fileConfig == "" {
		return nil, errs.ErrEmptyConfigPath
	}

	f, err := os.Open(fileConfig)
	if err != nil {
		return nil, errors.Join(errs.ErrFailedOpenConfig, err)
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)

	return read(f)
}

func read(reader io.Reader) (*Config, error) {
	if reader == nil {
		return nil, errs.ErrIncorrectReader
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errs.ErrFailedReadBuffer
	}

	var cfg Config

	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Join(errs.ErrFailedParseConfig, err)
	}

	return &cfg, nil
}
