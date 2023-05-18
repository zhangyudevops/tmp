package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type sConfig struct{}

func Config() *sConfig {
	return &sConfig{}
}

func (s *sConfig) ParseConfig(ctx context.Context, pattern string) (string, error) {
	cfg, err := g.Config().Get(ctx, pattern)
	if err != nil {
		return "", err
	}
	return cfg.String(), nil
}
