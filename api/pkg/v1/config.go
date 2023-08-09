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

type UpdateVarConfigReq struct {
	g.Meta `path:"/config/update" method:"post" tags:"更新配置"`
	Id     int64  `json:"id"`
	Value  string `json:"value"`
}

type UpdateVarConfigRes struct {
	Config string `json:"config"`
}
