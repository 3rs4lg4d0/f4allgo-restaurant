package test

import (
	"context"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/uber-go/tally/v4"
)

// nopTrManager is a transaction manager that does nothing (it just calls the
// function passed as argument and return the result).
type nopTrManager struct{}

var _ trm.Manager = (*nopTrManager)(nil)

func NewNopTrManager() *nopTrManager {
	return &nopTrManager{}
}

func (*nopTrManager) Do(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

func (*nopTrManager) DoWithSettings(ctx context.Context, _ trm.Settings, f func(ctx context.Context) error) error {
	return f(ctx)
}

// nopTallyTimer is a Tally timer that does nothing.
type nopTallyTimer struct{}

var _ tally.Timer = (*nopTallyTimer)(nil)
var _ tally.StopwatchRecorder = (*nopTallyTimer)(nil)

func NewNopTallyTimer() *nopTallyTimer {
	return &nopTallyTimer{}
}

func (*nopTallyTimer) Record(value time.Duration) {
	// NOP
}

func (ntt *nopTallyTimer) Start() tally.Stopwatch {
	return tally.NewStopwatch(time.Now(), ntt)
}

func (*nopTallyTimer) RecordStopwatch(stopwatchStart time.Time) {
	// NOP
}
