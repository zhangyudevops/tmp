package k8s

import (
	"context"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/model"
	"pack/internal/service"
)

func (c *cK8S) GetDeployPods(ctx context.Context, req *apiv1.ListPodsReq) (res *apiv1.ListPodsRes, err error) {
	ret, err := service.K8S().GetDeployPods(ctx, req.Namespace, req.Labels)
	if err != nil {
		return nil, err
	}

	res = &apiv1.ListPodsRes{
		List: ret,
	}

	return
}

func (c *cK8S) CopyFileToPod(ctx context.Context, req *apiv1.CopyFileToPodReq) (res *apiv1.CopyFileToPodRes, err error) {
	containerInfo := &model.Execute{
		Namespace:     req.Namespace,
		PodName:       req.Pod,
		ContainerName: req.Container,
	}

	if err = service.K8S().CopyFileToPod(ctx, containerInfo, req.InPath, req.OutPath); err != nil {
		return
	}

	return
}
