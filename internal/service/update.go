package service

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"pack/internal/dao"
	"pack/internal/model/entity"
	"strings"
)

type sUpdate struct{}

func Update() *sUpdate {
	return &sUpdate{}
}

// uncompressedUpdatePackage 解压升级压缩包，并返回解压后的目录
func (s *sUpdate) uncompressedUpdatePackage(ctx context.Context) (err error, path string) {
	pkgPath, _ := g.Config().Get(ctx, "package.path")
	packageDir := pkgPath.String() + "/tmp"

	// if the directory was existed, delete it
	if gfile.Exists(packageDir) {
		if err = gfile.Remove(packageDir); err != nil {
			return
		}
	}

	// get the newest package
	updatePkgPath, err := File().GetNewestDir(ctx, pkgPath.String())
	if err != nil {
		return
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
	imagesPath := path
	g.Log().Debugf(ctx, "imagesPath: %s", imagesPath)
	files, err := Path().GetFile(ctx, imagesPath, "*.tar")
	if err != nil {
		return
	}

	// 循环文件，导入并传仓库
	var loadImage, newImage string
	for _, file := range files {
		err, loadImage = Docker().LoadImage(ctx, file)
		if err != nil {
			return
		}

		// 修改tag
		err, newImage = Docker().TagDockerImage(ctx, loadImage)
		if err != nil {
			return
		}

		// 上传镜像到harbor
		if err = Docker().PushDockerImage(ctx, newImage); err != nil {
			return
		}

		// 把镜像信息写入数据库
		var imageInfo *entity.Image
		nameIndex := strings.LastIndex(newImage, "/")
		tagIndex := strings.Index(newImage, ":")
		name := newImage[nameIndex+1 : tagIndex]

		imageInfo = &entity.Image{
			Name: name,
			Tag:  newImage,
		}

		// 查询数据库中是否存在该条记录
		var ret gdb.Record
		ret, err = dao.Image.Ctx(ctx).Where("name", name).
			Where("tag", newImage).One()
		if err != nil {
			return err
		}
		if ret.IsEmpty() {
			// 不存在，写入
			if _, err = dao.Image.Ctx(ctx).Data(imageInfo).Insert(); err != nil {
				return err
			}
		}
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
