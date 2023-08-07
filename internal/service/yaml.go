package service

import (
	"context"
	"os"
	"pack/internal/dao"
	"pack/internal/model/entity"
	"pack/utility/util"
	"text/template"
)

type sYaml struct{}

func Yaml() *sYaml {
	return &sYaml{}
}

// RenderYamlFile 渲染yaml模版，输出到指定的目录
// 变量的key不能以.来区分，会识别不到，所以这里用_来区分，
// key的名称namespace_appName_key来命名
func (s *sYaml) RenderYamlFile(ctx context.Context, inFile, outFile string) (err error) {
	// 创建模版对象
	tmpl, err := template.ParseFiles(inFile)
	if err != nil {
		return
	}

	// 获取配置数据
	// 改了表结构，这里需要重构，把查询出的结果转变为map
	var (
		setting []*entity.Config
		config  = make(map[string]interface{})
	)
	err = dao.Config.Ctx(ctx).Scan(&setting)
	if err != nil {
		return
	}
	for _, set := range setting {
		// 把name和value的值转换为map对应的key和value
		config[set.Name] = set.Value
	}

	// 获取镜像列表
	v, err := dao.Image.Ctx(ctx).Fields("name,tag").Where("status=?", "1").All()
	if err != nil {
		return
	}

	imageEnv := make(map[string]interface{})
	for _, record := range v.List() {
		imageEnv[record["name"].(string)] = record["tag"].(string)
	}

	config = util.MergeMap(config, imageEnv)
	// 创建输出的yaml文件
	out, err := os.Create(outFile)
	if err != nil {
		return
	}
	defer out.Close()

	// 渲染模版写入文件
	if err = tmpl.Execute(out, config); err != nil {
		return
	}

	return
}
