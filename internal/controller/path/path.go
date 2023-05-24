package path

import (
	"context"
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
	list, err := service.Path().GetFile(ctx, req.Path, req.Pattern)
	if err != nil {
		return nil, err
	}

	res = &apiv1.FilesOrDirsListRes{
		Name: list,
	}

	return
}
