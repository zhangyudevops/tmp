package k8s

import (
	"context"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
)

// ListNamespaces 列出namespace
func (c *cK8S) ListNamespaces(ctx context.Context, req *apiv1.ListNamespaceReq) (res *apiv1.ListNamespaceRes, err error) {
	ret, err := service.K8S().GetNamespaces(ctx)
	if err != nil {
		return
	}
	res = &apiv1.ListNamespaceRes{
		List: ret,
	}
	return
}
