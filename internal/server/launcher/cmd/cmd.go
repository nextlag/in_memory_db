package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/nextlag/in_memory_db/internal/errs"
	"github.com/nextlag/in_memory_db/internal/server/entity"
)

const emptyRequest = 0

type Console struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Console {
	return &Console{log: log}
}

func (c *Console) ServerRun(ctx context.Context, handler entity.Handler) error {
	var response []byte

	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fmt.Print("> ")

		request, err := reader.ReadString('\n')
		if err != nil {
			c.log.Error("error read string", errs.ErrLog(err))

			continue

		}

		request = strings.TrimSpace(request)
		query := []byte(request)

		if len(request) == emptyRequest {
			continue
		}

		response = handler(ctx, query)

		switch string(response) {
		case entity.ResponseOk:
			color.Green(string(response))
		default:
			color.Red(string(response))
		}
	}
}

func (c *Console) Close() {
	c.log.Info("closing cmd launcher")
}
