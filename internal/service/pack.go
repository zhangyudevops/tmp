package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

type sPack struct{}

func Pack() *sPack {
	return &sPack{}
}

// PackTarFiles pack tar files
// if the pkg.path is /data/pkg, the newest pkg is /data/pkg/tmp
// check the directory /data/pkg/tmp, if the directory was existed, delete it, and create a new directory
func (s *sPack) PackTarFiles(ctx context.Context, images []string) (err error) {
	pkgPath, _ := g.Config().Get(ctx, "pkg.path")
	packageDir := pkgPath.String() + "/tmp"
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
