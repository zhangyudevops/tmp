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

// CompressTarGzip compress path directory to tar.gz file
// @todo: need to be tested
func (s *sFile) CompressTarGzip(ctx context.Context, path, name string) error {
	// create tar.gz file
	fw, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fw.Close()
	// gzip writer
	gw := gzip.NewWriter(fw)
	defer gw.Close()
	// tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()
	// get file list
	list, err := gfile.ScanDir(path, "*", true)
	if err != nil {
		return err
	}
	// write file
	for _, file := range list {
		// open file
		fr, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fr.Close()
		// get file info
		fi, err := fr.Stat()
		if err != nil {
			return err
		}
		// write file header
		hdr := &tar.Header{
			Name:    fi.Name(),
			Size:    fi.Size(),
			Mode:    int64(fi.Mode()),
			ModTime: fi.ModTime(),
		}
		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}
		// write file content
		if _, err = io.Copy(tw, fr); err != nil {
			return err
		}
	}
	return nil
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
				newPath = s2
			}

		}
	}

	return
}
