package service

import (
	"log/slog"
	"time"
)

type BatchCounter[K comparable] struct {
	events chan K
	ticker *time.Ticker
	flush  func(map[K]int) error
	done   chan struct{}
}

func NewBatchCounter[K comparable](bufSize int, interval time.Duration, flush func(map[K]int) error) *BatchCounter[K] {
	return &BatchCounter[K]{
		events: make(chan K, bufSize),
		ticker: time.NewTicker(interval),
		flush:  flush,
		done:   make(chan struct{}),
	}
}

func (b *BatchCounter[K]) Submit(key K) {
	select {
	case b.events <- key:
	default:
		slog.Warn("batch counter channel full, dropping event")
	}
}

func (b *BatchCounter[K]) Close() <-chan struct{} {
	close(b.events)
	return b.done
}

func (b *BatchCounter[K]) Run() {
	defer close(b.done)
	counts := make(map[K]int)
	defer b.ticker.Stop()

	flushNow := func() {
		if len(counts) == 0 {
			return
		}
		if err := b.flush(counts); err != nil {
			slog.Error("batch flush failed", "error", err)
		}
		counts = make(map[K]int)
	}

	for {
		select {
		case key, ok := <-b.events:
			if !ok {
				flushNow()
				return
			}
			counts[key]++
		case <-b.ticker.C:
			flushNow()
		}
	}
}
