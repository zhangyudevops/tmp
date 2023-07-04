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
		path, _ := g.Config().Get(ctx, "package.path")
		req.Path = path.String()
	}
	// if req.Path is directory, list the first layer files and dirs, if req.Path is file, return the file name
	var list []string
	if gfile.IsDir(req.Path) {
		g.Log().Debugf(ctx, "path %s is a directory", req.Path)
		list, err = service.Path().GetFile(ctx, req.Path, req.Pattern)
		if err != nil {
			return nil, err
		}
	} else {
		list = append(list, req.Path)
		g.Log().Debugf(ctx, "path %s is a file", req.Path)
	}

	res = &apiv1.FilesOrDirsListRes{
		Name: list,
	}

	return
}
