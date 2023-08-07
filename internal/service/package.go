package service

import (
	"context"
	"pack/internal/dao"
	"pack/internal/model/entity"
)

type sPackage struct{}

func Package() *sPackage {
	return &sPackage{}
}

func (s *sPackage) GetPackageHistory(ctx context.Context) ([]*entity.Upload, error) {
	var historyList []*entity.Upload
	err := dao.Upload.Ctx(ctx).Limit(3).OrderDesc("up_time").Scan(&historyList)
	if err != nil {
		return nil, err
	}

	return historyList, nil
}
