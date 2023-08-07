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
