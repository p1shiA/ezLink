package domain

import "context"

// TxManager coordinates multi-repository database transactions using Context
type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
