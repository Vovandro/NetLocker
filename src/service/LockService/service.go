package LockService

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/NetLocker/src/repository/LockRepository"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/paranoia/service"
	"gitlab.com/devpro_studio/go_utils/decode"
	"sync"
)

type Service struct {
	service.Mock
	cfg            ServiceConfig
	poolShard      []chan dataSource
	wg             sync.WaitGroup
	lockRepository LockRepository.ILockRepository
}

type dataSource struct {
	key        string
	timeout    int64
	res        chan bool
	isDeleteOp bool
}

type ServiceConfig struct {
	ShardCount int `yaml:"shard_count"`
}

func New(name string) *Service {
	return &Service{
		Mock: service.Mock{
			NamePkg: name,
		},
	}
}

func (t *Service) Init(app interfaces.IEngine, cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.cfg, "yaml", decode.DecoderStrongFoundDst)

	if err != nil {
		return err
	}

	t.lockRepository = app.GetModule(interfaces.ModuleRepository, "lock").(LockRepository.ILockRepository)

	if t.cfg.ShardCount <= 0 {
		t.cfg.ShardCount = 1
	}

	t.init()

	return nil
}

func (t *Service) init() {
	t.poolShard = make([]chan dataSource, t.cfg.ShardCount)
	for i := 0; i < t.cfg.ShardCount; i++ {
		t.poolShard[i] = make(chan dataSource, 100)
		go t.worker(i)
	}

	t.wg.Add(t.cfg.ShardCount)
}

func (t *Service) Stop() error {
	for i := 0; i < t.cfg.ShardCount; i++ {
		close(t.poolShard[i])
	}

	t.wg.Wait()

	return nil
}

func (t *Service) worker(num int) {
	for {
		select {
		case data, ok := <-t.poolShard[num]:
			if !ok {
				t.wg.Done()
				return
			}

			if data.isDeleteOp {
				err := t.lockRepository.Unlock(data.key)
				if err != nil {
					data.res <- false
				} else {
					data.res <- true
				}
			} else {
				data.res <- t.lockRepository.TryAndLock(data.key, data.timeout)
			}

			close(data.res)
		}
	}
}

func (t *Service) Lock(ctx context.Context, key string, timeout int64) error {
	shard := t.getShardKey(key)
	data := dataSource{
		key:     key,
		timeout: timeout,
		res:     make(chan bool),
	}
	t.poolShard[shard] <- data

	select {
	case res := <-data.res:
		if res {
			return nil
		}

	case <-ctx.Done():
		return ctx.Err()
	}

	return fmt.Errorf("timeout")
}

func (t *Service) Unlock(ctx context.Context, key string) error {
	shard := t.getShardKey(key)
	data := dataSource{
		key:        key,
		isDeleteOp: true,
		res:        make(chan bool),
	}
	t.poolShard[shard] <- data

	select {
	case res := <-data.res:
		if res {
			return nil
		}

	case <-ctx.Done():
		return ctx.Err()
	}

	return fmt.Errorf("timeout")
}

func (t *Service) getShardKey(key string) int {
	hash := 0
	for i := 0; i < len(key); i++ {
		hash = (31*hash + int(key[i])) % t.cfg.ShardCount
	}
	return hash
}
