package k8s

import (
	"context"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
)

// ListDeploys 列出指定namespace下的deployment
func (c *cK8S) ListDeploys(ctx context.Context, req *apiv1.ListDeployReq) (res *apiv1.ListDeployRes, err error) {
	ret, err := service.K8S().GetDeploys(ctx, req.Namespace)
	if err != nil {
		return
	}

	res = &apiv1.ListDeployRes{
		List: ret,
	}
	return
}
