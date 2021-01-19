package references

import "jvmgo/ch08/instructions/base"
import "jvmgo/ch08/rtda"
import "jvmgo/ch08/rtda/heap"

// newarray指令用来创建基本类型数组，包括boolean[]、byte[]、
// char[]、short[]、int[]、long[]、float[]和double[]8种。
const (
	//Array Type  atype
	AT_BOOLEAN = 4
	AT_CHAR    = 5
	AT_FLOAT   = 6
	AT_DOUBLE  = 7
	AT_BYTE    = 8
	AT_SHORT   = 9
	AT_INT     = 10
	AT_LONG    = 11
)

// Create new array
type NEW_ARRAY struct {
	atype uint8
}

func (self *NEW_ARRAY) FetchOperands(reader *base.BytecodeReader) {
	self.atype = reader.ReadUint8()
}
func (self *NEW_ARRAY) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	count := stack.PopInt()
	if count < 0 {
		panic("java.lang.NegativeArraySizeException")
	}

	// 根据atype值使用当前类的类加载器加载数组类，然后创建数组对象并推入操作数栈。
	classLoader := frame.Method().Class().Loader()
	arrClass := getPrimitiveArrayClass(classLoader, self.atype)
	arr := arrClass.NewArray(uint(count))
	stack.PushRef(arr)
}

func getPrimitiveArrayClass(loader *heap.ClassLoader, atype uint8) *heap.Class {
	switch atype {
	case AT_BOOLEAN:
		return loader.LoadClass("[Z")
	case AT_BYTE:
		return loader.LoadClass("[B")
	case AT_CHAR:
		return loader.LoadClass("[C")
	case AT_SHORT:
		return loader.LoadClass("[S")
	case AT_INT:
		return loader.LoadClass("[I")
	case AT_LONG:
		return loader.LoadClass("[J")
	case AT_FLOAT:
		return loader.LoadClass("[F")
	case AT_DOUBLE:
		return loader.LoadClass("[D")
	default:
		panic("Invalid atype!")
	}
}
