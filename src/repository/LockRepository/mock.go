package LockRepository

import (
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Unlock(key string) error {
	args := m.Called(key)
	return args.Error(0)
}
func (m *Mock) TryAndLock(key string, timeout int64) bool {
	args := m.Called(key, timeout)
	return args.Bool(0)
}
