package clean

import (
	"context"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
)

type cClean struct{}

func Clean() *cClean {
	return &cClean{}
}

func (c *cClean) HarborImageClean(ctx context.Context, req *apiv1.CleanHarborReq) (res apiv1.CleanHarborRes, err error) {
	err = service.Clean().HarborImageClean(ctx)
	if err != nil {
		return
	}

	return
}
