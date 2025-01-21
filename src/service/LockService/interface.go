package LockService

import "context"

type ILockService interface {
	Lock(ctx context.Context, key string, timeout int64) error
	Unlock(ctx context.Context, key string) error
}
