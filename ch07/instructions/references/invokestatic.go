package references

import "jvmgo/ch07/instructions/base"
import "jvmgo/ch07/rtda"
import "jvmgo/ch07/rtda/heap"

// 表示执行静态方法
// Invoke a class (static) method
type INVOKE_STATIC struct{ base.Index16Instruction }

func (self *INVOKE_STATIC) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	resolvedMethod := methodRef.ResolvedMethod()
	if !resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	class := resolvedMethod.Class()
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

	base.InvokeMethod(frame, resolvedMethod)
}
