package LockRepository

import (
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Unlock(key string, id string) error {
	args := m.Called(key, id)
	return args.Error(0)
}
func (m *Mock) TryAndLock(key string, id string, timeout int64) bool {
	args := m.Called(key, id, timeout)
	return args.Bool(0)
}
