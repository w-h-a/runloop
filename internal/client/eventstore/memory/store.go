package memory

import (
	"context"
	"sync"

	"github.com/w-h-a/runloop/internal/client/eventstore"
	"github.com/w-h-a/runloop/internal/domain/run"
)

type memoryEventStore struct {
	options eventstore.Options
	mtx     sync.RWMutex
	events  map[string][]run.Event
}

func (s *memoryEventStore) Append(_ context.Context, runID string, events ...run.Event) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.events[runID] = append(s.events[runID], events...)

	return nil
}

func (s *memoryEventStore) Load(_ context.Context, runID string) ([]run.Event, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	stored := s.events[runID]
	if len(stored) == 0 {
		return []run.Event{}, nil
	}

	cpy := make([]run.Event, len(stored))
	copy(cpy, stored)

	return cpy, nil
}

func (s *memoryEventStore) ListRunIDs(_ context.Context) ([]string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	runIDs := make([]string, 0, len(s.events))
	for runID := range s.events {
		runIDs = append(runIDs, runID)
	}

	return runIDs, nil
}

func NewEventStore(opts ...eventstore.Option) eventstore.EventStore {
	options := eventstore.NewOptions(opts...)

	s := &memoryEventStore{
		options: options,
		mtx:     sync.RWMutex{},
		events:  make(map[string][]run.Event),
	}

	return s

}
