package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
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

// UpdateVariableConfig 更新变量配置
func (s *sConfig) UpdateVariableConfig(ctx context.Context, id int64, value string) (err error) {
	_, err = dao.Config.Ctx(ctx).Data(g.Map{"value": value}).Where("id=?", id).Update()
	if err != nil {
		return
	}

	if _, err = g.Redis().Do(ctx, "DEL", "config"); err != nil {
		return
	}

	return nil
}
