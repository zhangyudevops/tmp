package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type FileUploadReq struct {
	// File is the file to upload
	g.Meta `path:"/file/upload" method:"post" mime:"multipart/form-data" tag:"file" summary:"File to upload"`
	File   *ghttp.UploadFile `json:"file" type:"file" dc:"File to upload"`
	MD5    string            `json:"md5" dc:"File md5 value"`
}

type FileUploadRes struct {
}
