package LockRepository

type ILockRepository interface {
	Unlock(key string, id string) error
	TryAndLock(key string, id string, timeout int64) bool
}
