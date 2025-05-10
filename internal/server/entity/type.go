package entity

import "context"

type Handler = func(ctx context.Context, query []byte) []byte
