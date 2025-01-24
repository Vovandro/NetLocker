package LockRepository

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/paranoia/repository"
	"gitlab.com/devpro_studio/Paranoia/pkg/cache/memory"
	"gitlab.com/devpro_studio/go_utils/decode"
	"math/rand"
	"strconv"
	"time"
)

type MemoryRepository struct {
	repository.Mock
	cache memory.IMemory
	cfg   RepositoryConfig
}

func NewMemory(name string) *MemoryRepository {
	return &MemoryRepository{
		Mock: repository.Mock{
			NamePkg: name,
		},
	}
}

func (t *MemoryRepository) Init(app interfaces.IEngine, cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.cfg, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.cfg.EnableDoubleCheck && t.cfg.TimeCheck <= 0 {
		t.cfg.TimeCheck = 1000
	}

	t.cache = app.GetPkg(interfaces.PkgCache, "primary").(memory.IMemory)

	return nil
}

func (t *MemoryRepository) Unlock(key string) error {
	return t.cache.Delete(context.Background(), key)
}

func (t *MemoryRepository) TryAndLock(key string, timeout int64) bool {
	if t.cache.Has(context.Background(), key) {
		return false
	}

	rnd := rand.Int63()
	rndStr := strconv.FormatInt(rnd, 10)

	err := t.cache.Set(context.Background(), key, rndStr, time.Second*time.Duration(timeout))
	if err != nil {
		return false
	}
	if t.cfg.EnableDoubleCheck {
		time.Sleep(time.Duration(t.cfg.TimeCheck))
		val, err := t.cache.Get(context.Background(), key)
		if err != nil {
			return false
		}
		if val != rndStr {
			return false
		}
	}

	return true
}
