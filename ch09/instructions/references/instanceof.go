package references

import "jvmgo/ch09/instructions/base"
import "jvmgo/ch09/rtda"
import "jvmgo/ch09/rtda/heap"

// instanceof指令判断对象是否是某个类的实例（或者对象的类是否实现了某个接口），并把结果推入操作数栈。

// Determine if object is of given type
type INSTANCE_OF struct{ base.Index16Instruction }

func (self *INSTANCE_OF) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	ref := stack.PopRef()
	if ref == nil {
		stack.PushInt(0)
		return
	}

	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolvedClass()
	if ref.IsInstanceOf(class) {
		stack.PushInt(1)
	} else {
		stack.PushInt(0)
	}
}
