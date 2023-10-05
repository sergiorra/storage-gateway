package context_wrapper

import (
	"context"

	"github.com/google/uuid"
)

type CorrelationIDKey string

const (
	correlationIDKey CorrelationIDKey = "CorrelationID"
)

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

func GetCorrelationID(ctx context.Context) string {
	val := ctx.Value(correlationIDKey)
	if val == nil {
		val = uuid.New().String()
	}

	return val.(string)
}
