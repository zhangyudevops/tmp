package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"pack/internal/model/entity"
)

type GetPackageHistoryReq struct {
	g.Meta `path:"/pkg/history" method:"get"  tag:"pkg" summary:"get package upload history"`
}

type GetPackageHistoryRes struct {
	HistoryList []*entity.Upload `json:"history_list" dc:"Package upload history list"`
}
