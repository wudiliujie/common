package rzip

import (
	"archive/zip"
	"fmt"
	"github.com/wudiliujie/common/progressbar"
	"github.com/wudiliujie/common/rpath"
	"io"
	"log"
	"os"
	"path/filepath"
)

//带进度条的压缩
func Zip(zipFile string, fileList []string) error {

	totalSize := int64(0)
	for _, v := range fileList {
		totalSize += rpath.GetFileSize(v)
	}
	// 创建 zip 包文件
	fw, err := os.Create(zipFile)
	if err != nil {
		log.Fatal()
	}
	defer fw.Close()
	pgb := progressbar.NewOptions64(totalSize, progressbar.OptionSetBytes64(totalSize))

	// 实例化新的 zip.Writer
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for _, fileName := range fileList {
		pgb.Describe(filepath.Base(fileName))
		fr, err := os.Open(fileName)
		if err != nil {
			return err
		}

		fi, err := fr.Stat()
		if err != nil {
			_ = fr.Close()
			return err
		}
		// 写入文件的头信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			_ = fr.Close()
			return err
		}
		fh.Method = zip.Deflate
		w, err := zw.CreateHeader(fh)
		if err != nil {
			_ = fr.Close()

			return err
		}
		// 写入文件内容
		//需要监控数量
		_, err = copyBuffer(w, fr, nil, pgb)
		if err != nil {
			_ = fr.Close()
			return err
		}
		_ = fr.Close()

	}
	fmt.Println()
	return nil
}
func copyBuffer(dst io.Writer, src io.Reader, buf []byte, pgb *progressbar.ProgressBar) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			_ = pgb.Add(nr)
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

func Unzip(zipFile string) error {
	zr, err := zip.OpenReader(zipFile)
	defer zr.Close()
	if err != nil {
		return err
	}

	for _, file := range zr.File {
		// 如果是目录，则创建目录
		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(file.Name, file.Mode()); err != nil {
				return err
			}
			continue
		}
		// 获取到 Reader
		fr, err := file.Open()
		if err != nil {
			return err
		}

		fw, err := os.OpenFile(file.Name, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
		fw.Close()
		fr.Close()
	}
	return nil
}
