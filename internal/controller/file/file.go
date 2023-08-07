package file

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
)

type cFile struct{}

func File() *cFile {
	return &cFile{}
}

func (c *cFile) Upload(ctx context.Context, req *apiv1.FileUploadReq) (res *apiv1.FileUploadRes, err error) {
	err = service.File().Upload(ctx, req.File, req.MD5)
	if err != nil {
		g.Log().Error(ctx, err.Error())
		return
	}
	g.Log().Infof(ctx, "上传文件成功: %s", req.File.Filename)

	return
}
