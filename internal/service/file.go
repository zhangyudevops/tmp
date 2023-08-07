package service

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"io"
	"os"
	"pack/internal/dao"
	"path/filepath"
	"strings"
	"time"
)

type sFile struct{}

func File() *sFile {
	return &sFile{}
}

func (s *sFile) Upload(ctx context.Context, inFile *ghttp.UploadFile, md5 string) (err error) {
	// 获取文件上传路径
	pkgVar, _ := g.Cfg().Get(ctx, "package.path")
	err = Path().CreateDir(ctx, pkgVar.String())
	if err != nil {
		return
	}
	//查询当前包名的md5是否存在
	record, err := dao.Upload.Ctx(ctx).Where("name", inFile.Filename).Where("md5", md5).All()
	if err != nil {
		return
	}
	if record.IsEmpty() {
		// 保存上传文件
		_, err = inFile.Save(pkgVar.String())
		if err != nil {
			return
		}
		// 保存上传记录
		_, err = dao.Upload.Ctx(ctx).Insert(g.Map{
			"name":    inFile.Filename,
			"md5":     md5,
			"path":    pkgVar.String(),
			"up_time": time.Now(),
		})
		if err != nil {
			return
		}
		return
	}

	// 如果不存在 ，返回错误
	return fmt.Errorf("文件已经存在，轻勿重复上传")
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
func (s *sFile) CompressTarGzip(ctx context.Context, source, target string) error {
	// if the source is a file, then compress it, if the source is a directory, then compress all files in it
	if gfile.IsFile(source) {
		return s.compressFile(ctx, source, target)
	}
	return s.compressDir(ctx, source, target)
}

func (s *sFile) compressFile(ctx context.Context, source, target string) error {
	// create target file
	fw, err := os.Create(target)
	if err != nil {
		return err
	}
	defer fw.Close()

	// create gzip writer
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// create tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// open source file
	fr, err := os.Open(source)
	if err != nil {
		return err
	}
	defer fr.Close()

	// get file info
	fi, err := fr.Stat()
	if err != nil {
		return err
	}

	// write tar header
	hdr := new(tar.Header)
	hdr.Name = fi.Name()
	hdr.Size = fi.Size()
	hdr.Mode = int64(fi.Mode())
	hdr.ModTime = fi.ModTime()

	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	// write file content
	if _, err := io.Copy(tw, fr); err != nil {
		return err
	}

	return nil
}

func (s *sFile) compressDir(ctx context.Context, source, target string) error {
	// create target file
	fw, err := os.Create(target)
	if err != nil {
		return err
	}
	defer fw.Close()

	// create gzip writer
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// create tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// walk path
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		// get file info
		hdr, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// set header name
		hdr.Name = strings.TrimPrefix(path, source)

		// write tar header
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// if not a dir, then write file content
		if !info.IsDir() {
			// open source file
			fr, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fr.Close()

			// write file content
			if _, err := io.Copy(tw, fr); err != nil {
				return err
			}
		}

		return nil
	})
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
