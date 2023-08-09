package config

import (
	"context"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
)

type cConfig struct{}

func Config() *cConfig {
	return &cConfig{}
}

func (c *cConfig) GetVarConfig(ctx context.Context, req *apiv1.VarConfigReq) (res *apiv1.VarConfigRes, err error) {
	ret, err := service.Config().GetVariableConfig(ctx)
	if err != nil {
		return
	}

	res = &apiv1.VarConfigRes{
		Config: ret,
	}

	return
}

// UpdateVarConfig 更新变量配置
func (c *cConfig) UpdateVarConfig(ctx context.Context, req *apiv1.UpdateVarConfigReq) (res *apiv1.UpdateVarConfigRes, err error) {
	err = service.Config().UpdateVariableConfig(ctx, req.Id, req.Value)
	if err != nil {
		return
	}

	res = &apiv1.UpdateVarConfigRes{
		Config: req.Value,
	}
	return
}
