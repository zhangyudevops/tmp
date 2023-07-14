package service

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/util/gconv"
	"os"
	"pack/internal/dao"
	"pack/internal/model/entity"
	"text/template"
)

type sYaml struct{}

func Yaml() *sYaml {
	return &sYaml{}
}

// RenderYamlFile 渲染yaml模版，输出到指定的目录
// 变量的key不能以.来区分，会识别不到
func (s *sYaml) RenderYamlFile(ctx context.Context, inFile, outFile string) (err error) {
	// 创建模版对象
	tmpl, err := template.ParseFiles(inFile)
	if err != nil {
		return
	}

	// 获取配置数据
	var setting entity.Setting
	err = dao.Setting.Ctx(ctx).Scan(&setting)
	if err != nil {
		return
	}
	configs := setting.Config

	// 字符串转换为json
	var data interface{}
	err = json.Unmarshal([]byte(configs), &data)
	if err != nil {
		return
	}
	config := gconv.Map(data)

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
