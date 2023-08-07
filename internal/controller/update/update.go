package update

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	apiv1 "pack/api/pkg/v1"
	"time"
)

type cUpdate struct{}

func Update() *cUpdate {
	return &cUpdate{}
}

// UpdateSystem 升级
func (c *cUpdate) UpdateSystem(ctx context.Context, req *apiv1.UpdateAppReq) (res *apiv1.UpdateAppRes, err error) {

	time.Sleep(5 * time.Second)

	time.Sleep(10 * time.Second)

	time.Sleep(4 * time.Second)

	time.Sleep(3 * time.Second)

	g.Dump("req", req)
	return
}
