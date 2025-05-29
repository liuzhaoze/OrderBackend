package domain

import "context"

type Locker interface {
	Lock(ctx context.Context, key, value string) (bool, error)
	Unlock(ctx context.Context, key, value string) (bool, error)
}
