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

	theNewestPath, err := service.File().GetNewestDir(ctx, filePath)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "The newest path is: %s", theNewestPath)
	if err = service.Path().CopyFileAndDir(theNewestPath, CurrentPackPath); err != nil {
		_ = service.File().DeleteCurrentDir(ctx, CurrentPackPath)
		return
	}

	dstPath := CurrentPackPath + "/images"
	// uncompressed the images.tar.gz file
	if err = service.File().ExtraTarGzip(ctx, CurrentPackPath+"/images.tar.gz", dstPath); err != nil {
		_ = service.File().DeleteCurrentDir(ctx, CurrentPackPath)
		return
	} else {
		// delete the images.tar.gz file
		if err = gfile.Remove(CurrentPackPath + "/images.tar.gz"); err != nil {
			return
		}
	}

	// @todo: 只是作为测试，如果解压了传统上面的images.ta.gz文件，需要把里层images目录下的文件移动到dstPath目录下,并删掉images目录
	if gfile.Exists(dstPath + "/images") {
		if err = gfile.Move(dstPath+"/images", dstPath); err != nil {
			//_ = service.File().DeleteCurrentDir(ctx, dstPath+"/images")
			return
		}
	}

	// request images list pull from harbor and save it to local
	if err = service.Docker().PullImageAndSaveToLocal(ctx, dstPath, req.Images); err != nil {
		_ = service.File().DeleteCurrentDir(ctx, CurrentPackPath)
		return nil, err
	}

	// compress the today's directory
	if err = service.File().CompressTarGzip(ctx, dstPath, CurrentPackPath+"/"+"images.tar.gz"); err != nil {
		_ = service.File().DeleteCurrentDir(ctx, CurrentPackPath)
		g.Log().Errorf(ctx, "Compress the today's directory failed: %s", err.Error())
		return nil, err
	} else {
		// delete the today's directory
		_ = service.File().DeleteCurrentDir(ctx, dstPath)
		g.Log().Infof(ctx, "Compress the today's directory %s successfully", CurrentPackPath+"/"+"images.tar.gz")
	}

	return &apiv1.PackUpdatePkgRes{}, nil
}
