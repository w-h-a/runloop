package eventstore

import (
	"context"

	"github.com/w-h-a/runloop/internal/domain/run"
)

type EventStore interface {
	Append(ctx context.Context, runID string, events ...run.Event) error
	Load(ctx context.Context, runID string) ([]run.Event, error)
	ListRunIDs(ctx context.Context) ([]string, error)
}
