package v1

import "github.com/gogf/gf/v2/frame/g"

type PackUpdatePkgReq struct {
	g.Meta `path:"/pack/full-pkg" method:"post"  tag:"pack" summary:"path update full pkg"`
	Images []string `json:"images"`
}

type PackUpdatePkgRes struct {
}
