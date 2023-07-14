package service

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"pack/internal/dao"
	"pack/internal/model/entity"
	"path/filepath"
	"strings"
)

type sUpdate struct{}

func Update() *sUpdate {
	return &sUpdate{}
}

// uncompressedUpdatePackage 解压升级压缩包，并返回解压后的目录
// pPath: 升级包存放的目录
// ePath: 解压后的目录
func (s *sUpdate) uncompressedUpdatePackage(ctx context.Context, pPath, ePath string) (err error) {
	// if the directory was existed, delete it
	if gfile.Exists(ePath) {
		if err = gfile.Remove(ePath); err != nil {
			return
		}
	}

	// get the newest package
	updatePkgPath, err := File().GetNewestDir(ctx, pPath)
	if err != nil {
		return
	}

	// create the directory
	if err = Path().CreateDir(ctx, ePath); err != nil {
		return
	}

	// 解压升级包到指定目录
	if err = File().ExtraTarGzip(ctx, updatePkgPath, ePath); err != nil {
		return
	}

	return nil
}

// pushNewImageToHarbor 把镜像文件推到新的harbor仓库
// path: 离线镜像.tar包存放的目录
func (s *sUpdate) pushNewImageToHarbor(ctx context.Context, path string) (err error) {
	files, err := Path().GetFile(ctx, path, "*.tar")
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

func (s *sUpdate) modifyConfigMap(ctx context.Context, tempPath, yamlPath string) (err error) {
	// 获取目录下所有模版文件
	fileList, err := Path().GetFile(ctx, tempPath, "*.tmpl")
	if err != nil {
		return
	}

	// 创建模版文件输出为yaml文件的目录
	if err = Path().CreateDir(ctx, yamlPath); err != nil {
		return
	}

	// 循环模版文件，替换模版文件中的变量
	for _, file := range fileList {
		// 组装yaml文件输出的绝对路径
		fileName := filepath.Base(file)
		outPath := yamlPath + "/" + strings.TrimSuffix(fileName, ".tmpl")
		if err = Yaml().RenderYamlFile(ctx, file, outPath); err != nil {
			return
		}
	}

	return
}

func (s *sUpdate) createOrUpdateFromYaml(ctx context.Context, path string) (err error) {
	yamlList, err := Path().GetFile(ctx, path, "*.yaml")
	if err != nil {
		return
	}

	for _, yaml := range yamlList {
		if err = K8S().CreateOrUpdateFromYamlFile(ctx, yaml); err != nil {
			return
		}
	}

	return
}

// Update 使用升级包进行升级.
// 升级目录结构：./images.tar.gz、./tmpl、./static. 模版文件以.yaml.tmpl结尾,yaml 文件以.yaml结尾
func (s *sUpdate) Update(ctx context.Context) (err error) {
	pkgVar, _ := g.Cfg().Get(ctx, "package.path")
	pkgPath := pkgVar.String()
	extraPath := pkgPath + "/tmp"

	// 1. 解压升级包
	if err = s.uncompressedUpdatePackage(ctx, pkgPath, extraPath); err != nil {
		return
	}

	// 2. 解压images.tar.gz
	imageDir := extraPath + "/images"
	if err = File().ExtraTarGzip(ctx, extraPath+"/images.tar.gz", imageDir); err != nil {
		return
	}

	// 3. 推送本地镜像到harbor
	if err = s.pushNewImageToHarbor(ctx, imageDir); err != nil {
		return
	}

	// 4. 更新静态文件

	// 5. 替换yaml模版为正式的yaml文件
	tmplPath := extraPath + "/tmpl"
	yamlPath := filepath.Dir(tmplPath) + "/yaml"
	if err = s.modifyConfigMap(ctx, tmplPath, yamlPath); err != nil {
		return
	}

	// 6. 使用yaml文件进行部署或者升级
	if err = s.createOrUpdateFromYaml(ctx, yamlPath); err != nil {
		return
	}

	return
}
