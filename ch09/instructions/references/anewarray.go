package references

import "jvmgo/ch09/instructions/base"
import "jvmgo/ch09/rtda"
import "jvmgo/ch09/rtda/heap"

// anewarray指令用来创建引用类型数组。
/**
anewarray指令也需要两个操作数。第一个操作数是uint16索
引，来自字节码。通过这个索引可以从当前类的运行时常量池中找
到一个类符号引用，解析这个符号引用就可以得到数组元素的类。
第二个操作数是数组长度，从操作数栈中弹出。
**/
// Create new array of reference
type ANEW_ARRAY struct{ base.Index16Instruction }

func (self *ANEW_ARRAY) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	componentClass := classRef.ResolvedClass()

	// if componentClass.InitializationNotStarted() {
	// 	thread := frame.Thread()
	// 	frame.SetNextPC(thread.PC()) // undo anewarray
	// 	thread.InitClass(componentClass)
	// 	return
	// }

	stack := frame.OperandStack()
	count := stack.PopInt()
	if count < 0 {
		panic("java.lang.NegativeArraySizeException")
	}

	arrClass := componentClass.ArrayClass()
	arr := arrClass.NewArray(uint(count))
	stack.PushRef(arr)
}
