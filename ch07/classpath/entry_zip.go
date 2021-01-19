package classpath

import "archive/zip"
import "errors"
import "io/ioutil"
import "path/filepath"

type ZipEntry struct {
	// 存放目录的绝对路径
	absPath string
	zipRC   *zip.ReadCloser
}

// 先把参数转换成绝对路径，如果转换过程中出现错误，则调用panic函数终止运行，否则创建ZipEntry实例并返回
func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return &ZipEntry{absPath, nil}
}

// 从Zip文件中读取class文件
// 首先打开ZIP文件，如果这一步出错的话，直接返回。然后遍历ZIP压缩包中的文件，看能否找到class文件。如果能找到，则打开class文件，把内容读取出来，并返回。
// 如果找不到，或者出现其他错误，则返回错误信息。有两处使用了defer语句来确保打开的文件得以关闭。
func (self *ZipEntry) readClass(className string) ([]byte, Entry, error) {
	if self.zipRC == nil {
		err := self.openJar()
		if err != nil {
			return nil, nil, err
		}
	}

	classFile := self.findClass(className)
	if classFile == nil {
		return nil, nil, errors.New("class not found: " + className)
	}

	data, err := readClass(classFile)
	return data, self, err
}

// todo: close zip
func (self *ZipEntry) openJar() error {
	r, err := zip.OpenReader(self.absPath)
	if err == nil {
		self.zipRC = r
	}
	return err
}

func (self *ZipEntry) findClass(className string) *zip.File {
	for _, f := range self.zipRC.File {
		if f.Name == className {
			return f
		}
	}
	return nil
}

func readClass(classFile *zip.File) ([]byte, error) {
	rc, err := classFile.Open()
	if err != nil {
		return nil, err
	}
	// read class data
	data, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 直接返回目录
func (self *ZipEntry) String() string {
	return self.absPath
}
