package references

import "jvmgo/ch06/instructions/base"
import "jvmgo/ch06/rtda"
import "jvmgo/ch06/rtda/heap"

// new指令专门用来创建类实例。数组由专门的指令创建，
// Create new object
type NEW struct{ base.Index16Instruction }

/**
因为接口和抽象类都不能实例化，所以如果解析后的类是接
口或抽象类，按照Java虚拟机规范规定，需要抛出InstantiationError
异常。另外，如果解析后的类还没有初始化，则需要先初始化类。
**/
func (self *NEW) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolvedClass()
	// todo: init class

	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationError")
	}

	ref := class.NewObject()
	frame.OperandStack().PushRef(ref)
}
