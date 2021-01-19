package references

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

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
	/**
	先判
	断类的初始化是否已经开始，如果还没有，则需要调用类的初始化
	方法，并终止指令执行。但是由于此时指令已经执行到了一半，也
	就是说当前帧的nextPC字段已经指向下一条指令，所以需要修改
	nextPC，让它重新指向当前指令。
	**/
	if !class.InitStarted() {
		frame.RevertNextPC()
		base.InitClass(frame.Thread(), class)
		return
	}

	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationError")
	}

	ref := class.NewObject()
	frame.OperandStack().PushRef(ref)
}
