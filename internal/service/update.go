package service

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"pack/internal/dao"
	"pack/internal/model"
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
			if _, err = dao.Image.Ctx(ctx).Where("name", name).Data(g.Map{"status": 0}).Update(); err != nil {
				return err
			}
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

// based directory create yaml file from template file,
// path is the tmpl config directory absolute path
func (s *sUpdate) modifyDirConfigFile(ctx context.Context, path string) (err error) {
	// modify service directory
	servicePath := path + "/service"
	if err = s.modifyConfigMap(ctx, servicePath, servicePath); err != nil {
		return
	}

	// modify configmap directory
	configmapPath := path + "/configmap"
	if err = s.modifyConfigMap(ctx, configmapPath, configmapPath); err != nil {
		return
	}

	// modify source directory
	sourcePath := path + "/source"
	if err = s.modifyConfigMap(ctx, sourcePath, sourcePath); err != nil {
		return
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

// based different directory, use yaml file create or update k8s resource,
// path is the tmpl directory absolute path
func (s *sUpdate) createOrUpdateFromYamlFile(ctx context.Context, path string) (err error) {
	// create or update service
	servicePath := path + "/service"
	if err = s.createOrUpdateFromYaml(ctx, servicePath); err != nil {
		return
	}

	// create or update configmap
	configmapPath := path + "/configmap"
	if err = s.createOrUpdateFromYaml(ctx, configmapPath); err != nil {
		return
	}

	// create or update source
	sourcePath := path + "/source"
	if err = s.createOrUpdateFromYaml(ctx, sourcePath); err != nil {
		return
	}

	return
}

// getCopyToPod 获取需要拷贝文件的pod名称,
// file 为要拷贝的文件或者目录的绝对路径,appName为应用名
func (s *sUpdate) getCopyToPodList(ctx context.Context, namespace, appName string) (err error, podList []*model.Pod) {
	// 获取label
	ret, err := K8S().DescribeDeploy(ctx, namespace, appName)
	if err != nil {
		return
	}
	label := ret.Labels

	// 获取pod
	podList, err = K8S().GetDeployPods(ctx, namespace, label)
	if err != nil {
		return
	}

	return
}

// copyFileToPod 把目录下的文件拷贝到pod中;
// file 如果为目录，则命名方式为 namespace_source_name,
// file 为要拷贝的文件或者目录的绝对路径
func (s *sUpdate) copyFileToPod(ctx context.Context, file string) (err error) {
	var (
		nameIndex          int
		namespace, appName string
	)
	dirName := filepath.Base(file)
	switch gfile.IsDir(file) {
	case true:
		// 如果是目录，截取base目录
		nameIndex = strings.Index(dirName, "_")
		namespace = dirName[:nameIndex]
		//source := dirName[nameIndex : nameIndex+1]
		appName = dirName[nameIndex+1:]

	case false:
		// 如果是文件，去掉扩展名
	}

	// 获取pod列表
	err, podList := s.getCopyToPodList(ctx, namespace, appName)
	if err != nil {
		return
	}

	// 获取pod名称
	for i, pod := range podList {
		if i == 0 {
			// 组装execute数据
			var execute *model.Execute
			execute = &model.Execute{
				Namespace: pod.Namespace,
				PodName:   pod.Name,
			}

			// 组装pod内部目标目录
			st, _ := dao.Static.Ctx(ctx).Where("name", appName).One()
			destPath := st.Map()["path"].(string)

			// 拷贝数据
			if err = K8S().CopyToToPod(ctx, execute, file, destPath); err != nil {
				return
			}
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

	// 4. 替换yaml模版为正式的yaml文件
	tmplPath := extraPath + "/tmpl"
	if err = s.modifyDirConfigFile(ctx, tmplPath); err != nil {
		return
	}

	// 5. 使用yaml文件进行部署或者升级
	if err = s.createOrUpdateFromYamlFile(ctx, tmplPath); err != nil {
		return
	}

	// 6. 更新静态文件
	staticDir := extraPath + "/static"
	dirList, err := gfile.ScanDir(staticDir, "*", false)
	if err != nil {
		return
	}
	for _, dir := range dirList {
		if err = s.copyFileToPod(ctx, dir); err != nil {
			return
		}
	}

	return
}
