package test

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"
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
