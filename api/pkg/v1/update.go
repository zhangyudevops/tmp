package v1

import "github.com/gogf/gf/v2/frame/g"

type UpdateAppReq struct {
	g.Meta `path:"/update/app" method:"post"  tag:"update" summary:"update app based selected package"`
	Name   string `json:"name" v:"required#name is required"`
}

type UpdateAppRes struct{}
