package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"pack/internal/model/entity"
)

type VarConfigReq struct {
	g.Meta `path:"/config/list" method:"get" tags:"列出配置详情"`
}

type VarConfigRes struct {
	Config []*entity.Config `json:"config"`
}
