package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CleanHarborReq struct {
	g.Meta `path:"/clean/harbor" method:"get"  tag:"clean" summary:"clean harbor images"`
}

type CleanHarborRes struct {
}

