package lang

import "fmt"
import "jvmgo/ch11/native"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

const jlThrowable = "java/lang/Throwable"

// 表示虚拟机栈帧信息
type StackTraceElement struct {
	// 帧正在执行的方法的的文件名
	fileName   string
	// 帧正在执行的方法的的类名
	className  string
	// 帧正在执行的方法的方法名
	methodName string
	// 帧正在执行哪行代码
	lineNumber int
}

func (self *StackTraceElement) String() string {
	return fmt.Sprintf("%s.%s(%s:%d)",
		self.className, self.methodName, self.fileName, self.lineNumber)
}

func init() {
	native.Register(jlThrowable, "fillInStackTrace", "(I)Ljava/lang/Throwable;", fillInStackTrace)
}

/**
查看java.lang.Exception或RuntimeException的源代码可以知道，
它们的构造函数都调用了超类java.lang.Throwable的构造函数。
Throwable的构造函数又调用了fillInStackTrace()方法记录Java虚拟机栈信息。
fillInStackTrace（）是用Java写的，必须借助另外一个本地方法
才能访问Java虚拟机栈，这个方法就是重载后的
fillInStackTrace（int）方法。
也就是说，要想抛出异常，Java虚拟机必须实现这个本地方法。
**/
// private native Throwable fillInStackTrace(int dummy);
// (I)Ljava/lang/Throwable;
func fillInStackTrace(frame *rtda.Frame) {
	this := frame.LocalVars().GetThis()
	frame.OperandStack().PushRef(this)

	stes := createStackTraceElements(this, frame.Thread())
	// 异常对象的extra字段中存放的就是Java虚拟机栈信息
	this.SetExtra(stes)
}

// 创建虚拟机栈帧信息
func createStackTraceElements(tObj *heap.Object, thread *rtda.Thread) []*StackTraceElement {
	/**
	由于栈顶两帧正在执行fillInStackTrace（int）和fillInStackTrace（）方法，
	所以需要跳过这两帧。这两帧下面的几帧正在执行异常类的构造函数，所以也要跳
	过，具体要跳过多少帧数则要看异常类的继承层次。
	**/
	skip := distanceToObject(tObj.Class()) + 2
	// 获取完整的Java虚拟机栈
	frames := thread.GetFrames()[skip:]
	stes := make([]*StackTraceElement, len(frames))
	for i, frame := range frames {
		stes[i] = createStackTraceElement(frame)
	}
	return stes
}

// 计算所需跳过的帧数
func distanceToObject(class *heap.Class) int {
	distance := 0
	for c := class.SuperClass(); c != nil; c = c.SuperClass() {
		distance++
	}
	return distance
}

// 根据帧创建虚拟机栈帧实例
func createStackTraceElement(frame *rtda.Frame) *StackTraceElement {
	method := frame.Method()
	class := method.Class()
	return &StackTraceElement{
		fileName:   class.SourceFile(),
		className:  class.JavaName(),
		methodName: method.Name(),
		lineNumber: method.GetLineNumber(frame.NextPC() - 1),
	}
}
