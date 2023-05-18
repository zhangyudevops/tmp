package pack

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
	"pack/utility/docker"
)

type cPack struct{}

func Pack() *cPack {
	return &cPack{}
}

func (c *cPack) PackUpdatePkg(ctx context.Context, req *apiv1.PackUpdatePkgReq) (res *apiv1.PackUpdatePkgRes, err error) {
	// create today's directory
	filePath, _ := service.Config().ParseConfig(ctx, "package.path")
	CurrentPackPath := filePath + "/" + docker.TodayDate()

	// if the current package directory was existed, delete it
	if gfile.Exists(CurrentPackPath) {
		err = gfile.Remove(CurrentPackPath)
		if err != nil {
			return nil, err
		}
	}

	// create the current package directory
	err = service.Path().CreateDir(ctx, CurrentPackPath)
	if err != nil {
		return nil, err
	}

	// copy the newest package directory to the current package directory
	scriptFilePath, _ := service.Config().ParseConfig(ctx, "script.path")
	sortFileShellScript := fmt.Sprintf("%s/list_dir_sorted.sh", scriptFilePath)
	theNewestPath, err := service.File().GetNewestPkgDir(ctx, sortFileShellScript, CurrentPackPath)
	if err != nil {
		return nil, err
	}
	if err = gfile.CopyDir(theNewestPath, CurrentPackPath); err != nil {
		return nil, err
	}

	// request images list pull from harbor and save it to local
	dstPath := CurrentPackPath + "/images"
	if err = service.Docker().PullImageAndSaveToLocal(ctx, dstPath, req.Images); err != nil {
		return nil, err
	}

	return &apiv1.PackUpdatePkgRes{}, nil
}
