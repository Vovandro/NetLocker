package NetLockerController

import (
	"context"
	"gitlab.com/devpro_studio/NetLocker/src/service/LockService"
	"gitlab.com/devpro_studio/Paranoia/paranoia/controller"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (t *Controller) TryAndLock(context.Context, *NetLockRequest) (*NetLockerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryAndLock not implemented")
}

func (t *Controller) Unlock(context.Context, *NetUnlockRequest) (*NetLockerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unlock not implemented")
}
