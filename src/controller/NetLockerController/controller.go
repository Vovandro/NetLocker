package NetLockerController

import (
	"context"
	"gitlab.com/devpro_studio/NetLocker/src/service/LockService"
	"gitlab.com/devpro_studio/Paranoia/paranoia/controller"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/pkg/server/grpc"
)

type Controller struct {
	lockService LockService.ILockService
	controller.Mock
	UnimplementedNetLockerServiceServer
}

func New(name string) *Controller {
	return &Controller{
		Mock: controller.Mock{
			NamePkg: name,
		},
	}
}

func (t *Controller) Init(app interfaces.IEngine, _ map[string]interface{}) error {
	t.lockService = app.GetModule(interfaces.ModuleService, "lock").(LockService.ILockService)

	app.GetPkg(interfaces.PkgServer, "grpc").(grpc.IGrpc).RegisterService(&NetLockerService_ServiceDesc, t)

	return nil
}

func (t *Controller) TryAndLock(c context.Context, req *NetLockRequest) (*NetLockerResponse, error) {
	resp := &NetLockerResponse{}

	err := t.lockService.Lock(c, req.Key, req.TimeLock)

	if err != nil {
		resp.Success = false
	} else {
		resp.Success = true
	}

	return resp, nil
}

func (t *Controller) Unlock(c context.Context, req *NetUnlockRequest) (*NetLockerResponse, error) {
	resp := &NetLockerResponse{}

	err := t.lockService.Unlock(c, req.Key)

	if err != nil {
		resp.Success = false
	} else {
		resp.Success = true
	}

	return resp, nil
}
