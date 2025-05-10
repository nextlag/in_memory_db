package common

import "context"

type TxID string

func ContextWithTxID(parent context.Context, value int64) context.Context {
	return context.WithValue(parent, TxID("tx"), value)
}

func GetTxIDFromContext(ctx context.Context) int64 {
	v := ctx.Value(TxID("tx"))

	id, ok := v.(int64)
	if !ok {
		return 0
	}

	return id
}
