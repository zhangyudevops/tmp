package docker

import (
	"context"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
)

type cDocker struct{}

func Docker() *cDocker {
	return &cDocker{}
}

func (c *cDocker) Test(ctx context.Context, req *apiv1.DockerTestReq) (res *apiv1.DockerTestRes, err error) {
	err = service.Docker().PushDockerImage(ctx, req.Image)
	if err != nil {
		return nil, err
	}

	return
}
