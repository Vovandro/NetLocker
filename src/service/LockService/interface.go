package LockService

import "context"

type ILockService interface {
	Lock(ctx context.Context, key string, id string, timeout int64) error
	Unlock(ctx context.Context, key string, id string) error
}
