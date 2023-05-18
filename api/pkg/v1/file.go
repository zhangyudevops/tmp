package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type FileUploadReq struct {
	// File is the file to upload
	g.Meta `path:"/file/upload" method:"post" mime:"multipart/form-data" tag:"file" summary:"File to upload"`
	File   *ghttp.UploadFiles `json:"file" type:"file" dc:"File to upload"`
	Path   string             `json:"path" dc:"File path"`
}

type FileUploadRes struct {
	Name []string `json:"name" dc:"File name"`
}
