package pkg

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
)

type cPkg struct{}

func Pkg() *cPkg {
	return &cPkg{}
}

// GetPackageHistory 获取上传包历史信息
func (c *cPkg) GetPackageHistory(ctx context.Context, req *apiv1.GetPackageHistoryReq) (res *apiv1.GetPackageHistoryRes, err error) {
	historyList, err := service.Package().GetPackageHistory(ctx)
	if err != nil {
		g.Log().Error(ctx, err.Error())
		return
	}
	res = &apiv1.GetPackageHistoryRes{
		HistoryList: historyList,
	}

	return
}
