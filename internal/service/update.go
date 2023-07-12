package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"pack/internal/dao"
)

type sUpdate struct{}

func Update() *sUpdate {
	return &sUpdate{}
}

// uncompressedUpdatePackage 解压升级压缩包，并返回解压后的目录
func (s *sUpdate) uncompressedUpdatePackage(ctx context.Context) (err error, path string) {
	pkgPath, _ := g.Config().Get(ctx, "package.path")
	packageDir := pkgPath.String() + "/tmp"

	// get the newest package
	updatePkgPath, err := File().GetNewestDir(ctx, pkgPath.String())
	if err != nil {
		return
	}

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

	// 解压升级包到指定目录, 现在不确定解压的目录层级
	if err = File().ExtraTarGzip(ctx, updatePkgPath, packageDir); err != nil {
		return
	}

	return nil, packageDir
}

func (s *sUpdate) PushNewImageToHarbor(ctx context.Context) (err error) {
	err, path := s.uncompressedUpdatePackage(ctx)
	if err != nil {
		return
	}

	// 获取images目录下面的所有镜像包
	imagesPath := path + "/images"
	files, err := Path().GetFile(ctx, imagesPath, "tar")
	if err != nil {
		return
	}

	// 循环文件，导入并传仓库
	for _, file := range files {
		if err = Docker().LoadImage(ctx, file); err != nil {
			return
		}

		// 获取镜像名称
	}

	return
}

func (s *sUpdate) getNewImageTag(ctx context.Context, name string) (error, string) {
	// 获取新镜像名称
	ret, err := dao.Image.Ctx(ctx).Where("name", name).OrderDesc("id").One()
	if err != nil {
		return err, ""
	}
	image := ret.Map()["tag"].(string)

	return nil, image
}
