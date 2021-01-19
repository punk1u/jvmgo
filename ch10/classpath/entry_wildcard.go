package classpath

import "os"
import "path/filepath"
import "strings"

// WildcardEntry实际上也是CompositeEntry
func newWildcardEntry(path string) CompositeEntry {
	baseDir := path[:len(path)-1] // remove *
	compositeEntry := []Entry{}

	// 根据后缀名选出JAR文件，并且返回SkipDir跳过子目录（通配符类路径不能递归匹配子目录下的JAR文件）
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != baseDir {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") {
			jarEntry := newZipEntry(path)
			compositeEntry = append(compositeEntry, jarEntry)
		}
		return nil
	}
	
	// 调用filepath包的Walk()函数遍历baseDir创建ZipEntry，第二个参数是上面定义好的函数
	filepath.Walk(baseDir, walkFn)

	return compositeEntry
}
