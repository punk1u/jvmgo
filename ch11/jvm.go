package main

import "fmt"
import "strings"
import "jvmgo/ch11/classpath"
import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

type JVM struct {
	cmd         *Cmd
	classLoader *heap.ClassLoader
	mainThread  *rtda.Thread
}

// 创建JVM结构体实例
func newJVM(cmd *Cmd) *JVM {
	cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)
	classLoader := heap.NewClassLoader(cp, cmd.verboseClassFlag)
	return &JVM{
		cmd:         cmd,
		classLoader: classLoader,
		mainThread:  rtda.NewThread(),
	}
}


func (self *JVM) start() {
	// 初始化VM类
	self.initVM()
	// 执行主类的main()方法
	self.execMain()
}

func (self *JVM) initVM() {
	// 加载VM类
	vmClass := self.classLoader.LoadClass("sun/misc/VM")
	// 执行VM类初始化方法
	base.InitClass(self.mainThread, vmClass)
	interpret(self.mainThread, self.cmd.verboseInstFlag)
}

// 先加载主类，然后执行其main()方法
func (self *JVM) execMain() {
	className := strings.Replace(self.cmd.class, ".", "/", -1)
	mainClass := self.classLoader.LoadClass(className)
	mainMethod := mainClass.GetMainMethod()
	if mainMethod == nil {
		fmt.Printf("Main method not found in class %s\n", self.cmd.class)
		return
	}

	// 在调用main()方法之前，需要给它传递args参数
	argsArr := self.createArgsArray()
	frame := self.mainThread.NewFrame(mainMethod)
	frame.LocalVars().SetRef(0, argsArr)
	self.mainThread.PushFrame(frame)
	interpret(self.mainThread, self.cmd.verboseInstFlag)
}

// 把Go的[]string变量转换成Java的字符串数组
func (self *JVM) createArgsArray() *heap.Object {
	stringClass := self.classLoader.LoadClass("java/lang/String")
	argsLen := uint(len(self.cmd.args))
	argsArr := stringClass.ArrayClass().NewArray(argsLen)
	jArgs := argsArr.Refs()
	for i, arg := range self.cmd.args {
		jArgs[i] = heap.JString(self.classLoader, arg)
	}
	return argsArr
}
