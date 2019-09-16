package rpath

import (
	"fmt"
	"github.com/wudiliujie/common/log"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if errcode != nil { //忽略错误
		// return errcode
		//}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}

		return nil
	})

	return files, err
}
func Mkdir(dirPath string, dirMode os.FileMode) error {
	dirPath = strings.Replace(dirPath, "/", "\\", -1)
	spath := strings.LastIndex(dirPath, "\\")
	dirPath = dirPath[0:spath]
	if FileExists(dirPath) {
		return nil
	}
	err := os.MkdirAll(dirPath, dirMode)
	if err != nil {
		return fmt.Errorf("%s: making directory: %v", dirPath, err)
	}
	return nil
}
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
func CheckAndCreatePath(name string) error {
	if !FileExists(name) {
		return Mkdir(name, 0777)
	}
	return nil
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}

	defer src.Close()
	err = Mkdir(dstName, 0777)
	if err != nil {
		log.Error("%v", err)
	}
	//路径不存在，直接创建路径
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func SaveFile(filename string, data string) {
	f, err1 := os.Create(filename) //创建文件
	if err1 != nil {
		fmt.Print(err1)
	}
	n, err1 := io.WriteString(f, data) //写入文件(字符串)
	if err1 != nil {
		fmt.Print(err1)
	}
	fmt.Printf("写入 %v 个字节\n", n)

}
func SaveBytes(filename string, data []byte) {
	f, err1 := os.Create(filename) //创建文件
	if err1 != nil {
		fmt.Print(err1)
		return
	}
	n, err1 := f.Write(data)
	if err1 != nil {
		fmt.Print(err1)
	}
	fmt.Printf("写入 %v 个字节\n", n)

}

func GetAllFile(pathname string) ([]string, error) {
	pathname = strings.TrimRight(pathname, "/")
	ret := make([]string, 0)
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return ret, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err := GetAllFile(fullDir)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return ret, err
			}
			ret = append(ret, s...)
		} else {
			fullName := pathname + "/" + fi.Name()
			ret = append(ret, fullName)
		}
	}
	return ret, nil
}
func GetFileSize(filename string) int64 {
	var result int64
	_ = filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
