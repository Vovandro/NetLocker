package LockRepository

type ILockRepository interface {
	Unlock(key string) error
	TryAndLock(key string, timeout int64) bool
}
