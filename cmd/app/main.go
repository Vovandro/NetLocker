package main

import (
	"context"
	"gitlab.com/devpro_studio/NetLocker/src/controller/NetLockerController"
	"gitlab.com/devpro_studio/NetLocker/src/repository/LockRepository"
	"gitlab.com/devpro_studio/NetLocker/src/service/LockService"
	"gitlab.com/devpro_studio/Paranoia/paranoia"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/pkg/cache/memory"
	"gitlab.com/devpro_studio/Paranoia/pkg/cache/redis"
	std_log "gitlab.com/devpro_studio/Paranoia/pkg/logger/std-log"
	"gitlab.com/devpro_studio/Paranoia/pkg/server/grpc"
	"gitlab.com/devpro_studio/Paranoia/pkg/server/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := paranoia.New("NetLocker", "cfg.yaml")

	s.PushPkg(std_log.New("std")).
		PushModule(LockService.New("lock")).
		PushModule(NetLockerController.New("controller"))

	cfg := s.GetConfig()

	if len(cfg.GetConfigItem(interfaces.PkgServer, "grpc")) > 0 {
		s.PushPkg(grpc.New("grpc"))
	}

	if len(cfg.GetConfigItem(interfaces.PkgServer, "http")) > 0 {
		s.PushPkg(http.New("http"))
	}

	switch cfg.GetString("cache_type", "redis") {
	case "redis":
		s.PushPkg(redis.New("primary")).
			PushModule(LockRepository.NewRedis("lock"))

	case "memory":
		s.PushPkg(memory.New("primary")).
			PushModule(LockRepository.NewMemory("lock"))
	}

	err := s.Init()

	if err != nil {
		panic(err)
	}

	defer s.Stop()

	s.GetLogger().Info(context.Background(), "start NetLocker service")

	// Wait for syscall stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
}
