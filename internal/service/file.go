package service

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
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

		out = fmt.Sprintf("%s/%s", dst, hdr.Name)
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

func (s *sFile) UnCompressFile(ctx context.Context, file, dst string) error {
	// 打开 .tgz 文件
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 创建 gzip.Reader
	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer gzipReader.Close()

	// 创建 tar.Reader
	tarReader := tar.NewReader(gzipReader)

	// 解压文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// 创建文件
		path := fmt.Sprintf("%s/%s", dst, header.Name)
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// 写入文件内容
		_, err = io.Copy(f, tarReader)
		if err != nil {
			panic(err)
		}
	}

	return err
}

// CompressTarGzip compress path directory to tar.gz file
func (s *sFile) CompressTarGzip(ctx context.Context, source, target string) error {
	// 创建目标文件
	targetFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	// 创建gzip压缩器
	gzipWriter := gzip.NewWriter(targetFile)
	defer gzipWriter.Close()

	// 创建tar打包器
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// 遍历源文件夹
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// 创建tar文件头
		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}

		// 写入文件头
		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件，写入文件内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (s *sFile) DeleteCurrentDir(ctx context.Context, dir string) error {
	return os.RemoveAll(dir)
}

func (s *sFile) GetNewestDir(ctx context.Context, pkgPath string) (newPath string, err error) {
	list, err := gfile.ScanDir(pkgPath, "*", false)
	if err != nil {
		return
	}
	if len(list) > 0 {
		// just get directory
		for i, s2 := range list {
			if !gfile.IsDir(s2) {
				continue
			}
			list[i] = strings.TrimSpace(s2)
		}

		// sort list by time
		var stat = time.Unix(0, 0).Unix()
		for _, s2 := range list {
			statPath, _ := gfile.Stat(s2)
			if stat < statPath.ModTime().Unix() {
				stat = statPath.ModTime().Unix()
				newPath = s2
			}

		}
	}

	return
}
