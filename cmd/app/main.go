package main

import (
	"context"
	"gitlab.com/devpro_studio/NetLocker/src/controller/NetLockerController"
	"gitlab.com/devpro_studio/NetLocker/src/repository/LockRepository"
	"gitlab.com/devpro_studio/NetLocker/src/service/LockService"
	"gitlab.com/devpro_studio/Paranoia/paranoia"
	"gitlab.com/devpro_studio/Paranoia/paranoia/telemetry"
	"gitlab.com/devpro_studio/Paranoia/pkg/cache/redis"
	std_log "gitlab.com/devpro_studio/Paranoia/pkg/logger/std-log"
	"gitlab.com/devpro_studio/Paranoia/pkg/server/grpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := paranoia.New("NetLocker", "cfg.yaml")

	s.SetMetrics(telemetry.NewPrometheusMetrics("metrics"))
	s.SetTrace(telemetry.NewTraceStd("trace"))

	s.PushPkg(std_log.New("std")).
		PushPkg(redis.New("primary")).
		PushPkg(grpc.New("grpc")).
		PushModule(LockRepository.New("lock")).
		PushModule(LockService.New("lock")).
		PushModule(NetLockerController.New("controller"))

	err := s.Init()

	if err != nil {
		panic(err)
	}

	defer s.Stop()

	s.GetLogger().Info(context.Background(), "start service")

	// Wait for syscall stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
}
