package pack

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
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
	if len(req.Images) == 0 {
		return nil, fmt.Errorf("images is empty")
	}
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
	//if err = service.Path().CreateDir(ctx, CurrentPackPath); err != nil {
	//	return
	//}

	// copy the newest package directory to the current package directory
	scriptFilePath, _ := service.Config().ParseConfig(ctx, "script.path")
	sortFileShellScript := fmt.Sprintf("%s/list_dir_sorted.sh", scriptFilePath)
	theNewestPath, err := service.File().GetNewestPkgDir(ctx, sortFileShellScript, filePath)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	g.Log().Debugf(ctx, "The newest path is: %s", theNewestPath)
	// if the newest path is empty, return error
	if theNewestPath == "" {
		return nil, fmt.Errorf("the newest path is empty")
	}
	// copy the newest dir to the current update dir
	if err = service.Path().CopyFileAndDir(theNewestPath, CurrentPackPath); err != nil {
		_ = service.File().DeleteCurrentDir(ctx, CurrentPackPath)
		return
	}

	// uncompressed the images.tar.gz file
	if err = service.File().ExtraTarGzip(ctx, CurrentPackPath+"/images.tar.gz", CurrentPackPath); err != nil {
		_ = service.File().DeleteCurrentDir(ctx, CurrentPackPath)
		return
	} else {
		// delete the images.tar.gz file
		if err = gfile.Remove(CurrentPackPath + "/images.tar.gz"); err != nil {
			return
		}
	}

	// request images list pull from harbor and save it to local
	dstPath := CurrentPackPath + "/images"
	if err = service.Docker().PullImageAndSaveToLocal(ctx, dstPath, req.Images); err != nil {
		_ = service.File().DeleteCurrentDir(ctx, CurrentPackPath)
		return nil, err
	}

	return &apiv1.PackUpdatePkgRes{}, nil
}
