package service

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"io"
	"os"
	"strings"
)

type sFile struct{}

func File() *sFile {
	return &sFile{}
}

func (s *sFile) Upload(ctx context.Context, inFile *ghttp.UploadFiles, path string) (fileList []string, err error) {
	err = Path().CreateDir(ctx, path)
	if err != nil {
		return nil, err
	}
	fileList, err = inFile.Save(path)
	if err != nil {
		return nil, err
	}
	return fileList, nil
}

// ExtraTarGzip 解压tar.gz文件
// file: 文件路径
// dst: 解压后的文件存放路径
// return: err, 解压后的文件路径
func (s *sFile) ExtraTarGzip(ctx context.Context, file, dst string) error {
	var out string
	// 读取文件
	fr, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fr.Close()
	// 解压
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	// 遍历文件
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		out = fmt.Sprintf("%s%s", dst, hdr.Name)
		// 判断文件类型
		switch hdr.Typeflag {
		case tar.TypeDir:
			// 创建文件夹
			if err = os.MkdirAll(out, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// 创建文件
			fw, err := os.Create(out)
			if err != nil {
				return err
			}
			// 写入文件
			if _, err = io.Copy(fw, tr); err != nil {
				return err
			}
			fw.Close()
		}
	}

	return err
}

// GetNewestPkgDir 获取pkg目录下最新的包目录
func (s *sFile) GetNewestPkgDir(ctx context.Context, file, pkgPath string) (newPath string, err error) {
	// 获取pkg目录下的所有目录
	argsList := make([]string, 1)
	argsList[0] = pkgPath
	g.Log().Debugf(ctx, "The script is %s %s", file, argsList)
	err, bytes := Shell().Exec(ctx, file, argsList)
	if err != nil {
		g.Log().Error(ctx, err)
		return "", err
	}

	// through the lines in the output, join to a slice of strings
	pkgDirList := strings.Split(string(bytes), "\n")
	for i, s2 := range pkgDirList {
		if strings.Contains(s2, " ") {
			continue
		}
		pkgDirList[i] = strings.TrimSpace(s2)

	}
	g.Log().Debugf(ctx, "Under the directory %s has: %s", pkgPath, pkgDirList)
	// 获取最新的包目录
	if len(pkgDirList) > 0 {
		g.Log().Debugf(ctx, "The newest pkg dir is %s", pkgDirList[len(pkgDirList)-1])
		if pkgDirList[len(pkgDirList)-1] == " " {
			newPath = pkgDirList[len(pkgDirList)-2]
		}
		newPath = pkgDirList[len(pkgDirList)-1]
	}
	return
}

func (s *sFile) DeleteCurrentDir(ctx context.Context, dir string) error {
	return os.RemoveAll(dir)
}
