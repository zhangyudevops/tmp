package service

import (
	"context"
	"github.com/gogf/gf/v2/os/gfile"
	"pack/utility/docker"
)

type sPack struct{}

func Pack() *sPack {
	return &sPack{}
}

// PackFullUpdateFile pack full update file
// if the package.path is /data/package, the newest package is /data/package/20210910
// if this function return error, the current package directory will be deleted
func (s *sPack) PackFullUpdateFile(ctx context.Context, images []string) (err error) {
	// find the newest path under pkg directory
	pkgPath, _ := Config().ParseConfig(ctx, "package.path")
	newestPath, err := File().GetNewestPkgDir(ctx, pkgPath, pkgPath)
	if err != nil {
		return
	}

	// create the current update dir
	packageDir := pkgPath + "/" + docker.TodayDate()
	if err = Path().CreateDir(ctx, packageDir); err != nil {
		return
	}

	// copy the newest dir to the current update dir
	if err = Path().copyFileAndDir(newestPath, packageDir); err != nil {
		_ = File().deleteCurrentDir(ctx, packageDir)
		return
	}

	// uncompressed the images.tar.gz file
	dstPath := packageDir + "/images"
	if err = File().extraTarGzip(ctx, packageDir+"/images.tar.gz", dstPath); err != nil {
		_ = File().deleteCurrentDir(ctx, packageDir)
		return err
	}

	// pull the images and save to the images directory
	if err = Docker().PullImageAndSaveToLocal(ctx, dstPath, images); err != nil {
		_ = File().deleteCurrentDir(ctx, packageDir)
		return
	}

	return
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
		_ = File().deleteCurrentDir(ctx, packageDir)
		return
	}

	return
}
