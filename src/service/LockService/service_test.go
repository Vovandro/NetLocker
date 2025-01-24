package LockService

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/devpro_studio/NetLocker/src/repository/LockRepository"
	"testing"
	"time"
)

func TestService_Lock_Success(t *testing.T) {
	mockRepo := new(LockRepository.Mock)
	mockRepo.On("TryAndLock", "test-key", mock.Anything).Return(true)

	service := New("test-service")
	service.cfg.ShardCount = 2
	service.lockRepository = mockRepo
	service.init()
	defer service.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := service.Lock(ctx, "test-key", 1000)
	assert.NoError(t, err)
}

func TestService_Lock_Failure(t *testing.T) {
	mockRepo := new(LockRepository.Mock)
	mockRepo.On("TryAndLock", "test-key", mock.Anything).Return(false)

	service := New("test-service")
	service.cfg.ShardCount = 2
	service.lockRepository = mockRepo
	service.init()
	defer service.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := service.Lock(ctx, "test-key", 1000)
	assert.Error(t, err)
	assert.Equal(t, "timeout", err.Error())
}

func TestService_Unlock_Success(t *testing.T) {
	mockRepo := new(LockRepository.Mock)
	mockRepo.On("Unlock", "test-key").Return(nil)

	service := New("test-service")
	service.cfg.ShardCount = 2
	service.lockRepository = mockRepo
	service.init()
	defer service.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := service.Unlock(ctx, "test-key")
	assert.NoError(t, err)
}

func TestService_Unlock_Failure(t *testing.T) {
	mockRepo := new(LockRepository.Mock)
	mockRepo.On("Unlock", "test-key").Return(errors.New("mock error"))

	service := New("test-service")
	service.cfg.ShardCount = 2
	service.lockRepository = mockRepo
	service.init()
	defer service.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := service.Unlock(ctx, "test-key")
	assert.Error(t, err)
}

func TestService_getShardKey(t *testing.T) {
	service := New("test-service")
	service.cfg.ShardCount = 2
	service.init()

	shard := service.getShardKey("test-key")
	assert.GreaterOrEqual(t, shard, 0)
	assert.Less(t, shard, 3)
}
