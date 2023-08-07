package service

import (
	"context"
	"pack/internal/dao"
	"pack/internal/model/entity"
)

type sConfig struct{}

func Config() *sConfig {
	return &sConfig{}
}

func (s *sConfig) GetVariableConfig(ctx context.Context) ([]*entity.Config, error) {
	var setting []*entity.Config
	err := dao.Config.Ctx(ctx).Order("app").Scan(&setting)
	if err != nil {
		return nil, err
	}

	return setting, nil
}
