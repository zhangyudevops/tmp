package path

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
)

type cPath struct{}

func Path() *cPath {
	return &cPath{}
}

func (c *cPath) ListFilesOrDirs(ctx context.Context, req *apiv1.FilesOrDirsListReq) (res *apiv1.FilesOrDirsListRes, err error) {
	if req.Path == "" {
		path, _ := service.Config().ParseConfig(ctx, "package.path")
		req.Path = path
	}
	// 判断path是否是目录，如果是目录，就返回目录下的文件列表，如果是文件，就返回空
	var list []string
	if gfile.IsDir(req.Path) {
		g.Log().Debugf(ctx, "path %s is a directory", req.Path)
		list, err = service.Path().GetFile(ctx, req.Path, req.Pattern)
		if err != nil {
			return nil, err
		}
	} else {
		g.Log().Debugf(ctx, "path %s is a file", req.Path)
	}

	res = &apiv1.FilesOrDirsListRes{
		Name: list,
	}

	return
}
