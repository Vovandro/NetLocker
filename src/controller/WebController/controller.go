package WebController

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/devpro_studio/NetLocker/src/service/LockService"
	"gitlab.com/devpro_studio/Paranoia/paranoia/controller"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/pkg/server/http"
	"strconv"
)

type Controller struct {
	lockService LockService.ILockService
	controller.Mock
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

	srv := app.GetPkg(interfaces.PkgServer, "http").(http.IHttp)

	srv.PushRoute("GET", "/lock", t.TryAndLock, nil)
	srv.PushRoute("GET", "/unlock", t.Unlock, nil)

	return nil
}

func (t *Controller) TryAndLock(c context.Context, ctx http.ICtx) {
	key := ctx.GetRequest().GetQuery().Get("key")
	id := ctx.GetRequest().GetQuery().Get("unique_id")
	timeLockStr := ctx.GetRequest().GetQuery().Get("time_lock")
	timeLock, err := strconv.ParseInt(timeLockStr, 10, 64)
	if err != nil || timeLock <= 0 || key == "" {
		ctx.GetResponse().SetBody([]byte("invalid request data"))
		ctx.GetResponse().SetStatus(422)
		return
	}

	if id == "" {
		id = uuid.New().String()
	}

	locked := t.lockService.Lock(c, key, id, timeLock) == nil

	ctx.GetResponse().SetBody([]byte(strconv.FormatBool(locked)))
	ctx.GetResponse().SetStatus(200)
	ctx.GetResponse().Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func (t *Controller) Unlock(c context.Context, ctx http.ICtx) {
	key := ctx.GetRequest().GetQuery().Get("key")
	id := ctx.GetRequest().GetQuery().Get("unique_id")

	if key == "" {
		ctx.GetResponse().SetBody([]byte("invalid request data"))
		ctx.GetResponse().SetStatus(422)
		return
	}

	locked := t.lockService.Unlock(c, key, id) == nil

	ctx.GetResponse().SetBody([]byte(strconv.FormatBool(locked)))
	ctx.GetResponse().SetStatus(200)
	ctx.GetResponse().Header().Set("Content-Type", "text/plain; charset=utf-8")
}
