package classpath

import "os"
import "strings"

// 存放路径分隔符的常量
// :(linux/unix) or ;(windows)
const pathListSeparator = string(os.PathListSeparator)

// Entry接口有四个实现：DirEntry、ZipEntry、CompositeEntry和WildcardEntry
type Entry interface {
	// 寻找和加载class文件，参数是class文件的相对路径，路径用/分隔，文件名有.class后缀，比如java/lang/Object.class，
	// 返回值是读取到的字节数据、最终定位到class文件的Entry、以及错误信息,Go的函数或方法允许返回多个值，按照惯例，使用最后一个返回值作为错误信息
	// className: fully/qualified/ClassName.class
	readClass(className string) ([]byte, Entry, error)
	String() string
}

// 根据参数创建不同类型的Entry实例
func newEntry(path string) Entry {
	if strings.Contains(path, pathListSeparator) {
		return newCompositeEntry(path)
	}

	if strings.HasSuffix(path, "*") {
		return newWildcardEntry(path)
	}

	if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".ZIP") {

		return newZipEntry(path)
	}

	return newDirEntry(path)
}
