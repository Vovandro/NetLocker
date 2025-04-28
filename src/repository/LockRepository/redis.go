package LockRepository

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/paranoia/repository"
	"gitlab.com/devpro_studio/Paranoia/pkg/cache/redis"
	"gitlab.com/devpro_studio/go_utils/decode"
	"time"
)

type RedisRepository struct {
	repository.Mock
	cache redis.IRedis
	cfg   RepositoryConfig
}

type RepositoryConfig struct {
	EnableDoubleCheck bool          `yaml:"enable_double_check"`
	TimeCheck         time.Duration `yaml:"time_check"`
}

func NewRedis(name string) *RedisRepository {
	return &RedisRepository{
		Mock: repository.Mock{
			NamePkg: name,
		},
	}
}

func (t *RedisRepository) Init(app interfaces.IEngine, cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.cfg, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.cfg.EnableDoubleCheck && t.cfg.TimeCheck <= 0 {
		t.cfg.TimeCheck = time.Second
	}

	t.cache = app.GetPkg(interfaces.PkgCache, "primary").(redis.IRedis)

	return nil
}

func (t *RedisRepository) Unlock(key string, id string) error {
	if id != "" {
		oldId, err := t.cache.Get(context.Background(), key)
		if err == nil {
			if oldId != id {
				return fmt.Errorf("invalid unlock id")
			}
		} else {
			return nil
		}
	}

	return t.cache.Delete(context.Background(), key)
}

func (t *RedisRepository) TryAndLock(key string, id string, timeout int64) bool {
	oldId, err := t.cache.Get(context.Background(), key)
	if err == nil {
		return oldId == id
	}

	err = t.cache.Set(context.Background(), key, id, time.Second*time.Duration(timeout))
	if err != nil {
		return false
	}
	if t.cfg.EnableDoubleCheck {
		time.Sleep(t.cfg.TimeCheck)
		val, err := t.cache.Get(context.Background(), key)
		if err != nil {
			return false
		}
		if val != id {
			return false
		}
	}

	return true
}
