package main

import "flag"
import "fmt"
import "os"

// java [-options] class [args...]
type Cmd struct {
	helpFlag         bool
	versionFlag      bool
	// 控制是否把类加载信息输出到控制台
	verboseClassFlag bool
	// 控制是否把指令执行信息输出到控制台
	verboseInstFlag  bool
	cpOption         string
	// 指定jre目录的位置以使用JDK的启动类路径(BootStrap Classpath)来寻找和加载JAVA标准库中的类
	XjreOption       string
	class            string
	args             []string
}

func parseCmd() *Cmd {
	cmd := &Cmd{}

	flag.Usage = printUsage
	flag.BoolVar(&cmd.helpFlag, "help", false, "print help message")
	flag.BoolVar(&cmd.helpFlag, "?", false, "print help message")
	flag.BoolVar(&cmd.versionFlag, "version", false, "print version and exit")
	flag.BoolVar(&cmd.verboseClassFlag, "verbose", false, "enable verbose output")
	flag.BoolVar(&cmd.verboseClassFlag, "verbose:class", false, "enable verbose output")
	flag.BoolVar(&cmd.verboseInstFlag, "verbose:inst", false, "enable verbose output")
	flag.StringVar(&cmd.cpOption, "classpath", "", "classpath")
	flag.StringVar(&cmd.cpOption, "cp", "", "classpath")
	// 增加处理指定jre目录位置的命令行参数
	flag.StringVar(&cmd.XjreOption, "Xjre", "", "path to jre")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		cmd.class = args[0]
		cmd.args = args[1:]
	}

	return cmd
}

func printUsage() {
	fmt.Printf("Usage: %s [-options] class [args...]\n", os.Args[0])
	//flag.PrintDefaults()
}
