package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type FilesOrDirsListReq struct {
	g.Meta  `path:"/path/list" method:"get"  tag:"list" summary:"path list"`
	Path    string `json:"path" dc:"File path"`
	Pattern string `json:"pattern"`
}

type FilesOrDirsListRes struct {
	Name []string `json:"name" dc:"File name or dir name"`
}
