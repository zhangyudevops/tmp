package service

import (
	"context"
	"github.com/gogf/gf/v2/os/gfile"
)

type sPack struct{}

func Pack() *sPack {
	return &sPack{}
}

// PackTarFiles pack tar files
// if the package.path is /data/package, the newest package is /data/package/tmp
// check the directory /data/package/tmp, if the directory was existed, delete it, and create a new directory
func (s *sPack) PackTarFiles(ctx context.Context, images []string) (err error) {
	pkgPath, _ := Config().ParseConfig(ctx, "package.path")
	packageDir := pkgPath + "/tmp"
	// if the directory was existed, delete it
	if gfile.Exists(packageDir) {
		if err = gfile.Remove(packageDir); err != nil {
			return
		}
	}

	// create the directory
	if err = Path().CreateDir(ctx, packageDir); err != nil {
		return
	}

	// pull the images and save to the images directory
	if err = Docker().PullImageAndSaveToLocal(ctx, packageDir, images); err != nil {
		_ = File().DeleteCurrentDir(ctx, packageDir)
		return
	}

	return
}
