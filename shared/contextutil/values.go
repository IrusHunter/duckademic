package contextutil

import (
	"context"

	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/google/uuid"
)

// ContextKey is a type used as key for storing values in context.
// It is defined as a string but using a distinct type avoids collisions.
type ContextKey string

const (
	traceIDKey ContextKey = "traceID"
	claimsKey  ContextKey = "accessClaims"
)

// SetTraceID creates a new UUID and stores it in the context under traceIDKey.
// Returns a new context containing the trace ID.
func SetTraceID(ctx context.Context) context.Context {
	id := uuid.New().String()
	return context.WithValue(ctx, traceIDKey, id)
}

// GetTraceID retrieves the trace ID from the context.
// Panics if the context does not contain a valid trace ID.
func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		panic("failed to get traceID from context")
	}
	return traceID
}

// SetAccessClaims stores JWT claims in context
func SetAccessClaims(ctx context.Context, claims *jsonutil.AccessClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetAccessClaims retrieves JWT claims from context
func GetAccessClaims(ctx context.Context) *jsonutil.AccessClaims {
	claims, ok := ctx.Value(claimsKey).(*jsonutil.AccessClaims)
	if !ok {
		return nil
	}
	return claims
}
