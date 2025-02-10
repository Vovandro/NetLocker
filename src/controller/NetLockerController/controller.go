package NetLockerController

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

	if req.Key == "" {
		return nil, fmt.Errorf("invalid request data")
	}

	if req.UniqueId == nil || *req.UniqueId == "" {
		req.UniqueId = new(string)
		*req.UniqueId = uuid.New().String()
	}

	resp.Success = t.lockService.Lock(c, req.Key, *req.UniqueId, req.TimeLock) == nil

	return resp, nil
}

func (t *Controller) Unlock(c context.Context, req *NetUnlockRequest) (*NetLockerResponse, error) {
	resp := &NetLockerResponse{}

	if req.Key == "" {
		return nil, fmt.Errorf("invalid request data")
	}

	if req.UniqueId == nil {
		req.UniqueId = new(string)
	}

	resp.Success = t.lockService.Unlock(c, req.Key, *req.UniqueId) == nil

	return resp, nil
}
