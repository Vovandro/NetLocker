package LockRepository

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/paranoia/repository"
	"gitlab.com/devpro_studio/Paranoia/pkg/cache/redis"
	"math/rand"
	"strconv"
	"time"
)

type Repository struct {
	repository.Mock
	cache redis.IRedis
}

func New(name string) *Repository {
	return &Repository{
		Mock: repository.Mock{
			NamePkg: name,
		},
	}
}

func (t *Repository) Init(app interfaces.IEngine, cfg map[string]interface{}) error {
	t.cache = app.GetPkg(interfaces.PkgCache, "primary").(redis.IRedis)

	return nil
}

func (t *Repository) Unlock(key string) error {
	return t.cache.Delete(context.Background(), key)
}

func (t *Repository) TryAndLock(key string, timeout int64) bool {
	if t.cache.Has(context.Background(), key) {
		return false
	}

	rnd := rand.Int63()
	rndStr := strconv.FormatInt(rnd, 10)

	err := t.cache.Set(context.Background(), key, rndStr, time.Second*time.Duration(timeout))
	if err != nil {
		return false
	}
	time.Sleep(time.Millisecond)
	val, err := t.cache.Get(context.Background(), key)
	if err != nil {
		return false
	}
	if val != rndStr {
		return false
	}

	return true
}
