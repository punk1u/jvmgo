package comparisons

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"

// dcmpg和dcmpl指令用于比较double变量,类似fcmpg、fcmpl指令。
// 这两条指令和lcmp指令很像，但是除了比较的变量类型不同以
// 外，还有一个重要的区别。由于浮点数计算有可能产生NaN（Not a
// Number）值，所以比较两个浮点数时，除了大于、等于、小于之外，
// 还有第4种结果：无法比较。dcmpg和dcmpl指令的区别就在于对第4
// 种结果的定义。

// 也就是说，当两个double变量中至少有一个是NaN时，用dcmpg指
// 令比较的结果是1，而用dcmpl指令比较的结果是-1。

// Compare double
type DCMPG struct{ base.NoOperandsInstruction }

func (self *DCMPG) Execute(frame *rtda.Frame) {
	_dcmp(frame, true)
}

type DCMPL struct{ base.NoOperandsInstruction }

func (self *DCMPL) Execute(frame *rtda.Frame) {
	_dcmp(frame, false)
}

func _dcmp(frame *rtda.Frame, gFlag bool) {
	stack := frame.OperandStack()
	v2 := stack.PopDouble()
	v1 := stack.PopDouble()
	if v1 > v2 {
		stack.PushInt(1)
	} else if v1 == v2 {
		stack.PushInt(0)
	} else if v1 < v2 {
		stack.PushInt(-1)
	} else if gFlag {
		stack.PushInt(1)
	} else {
		stack.PushInt(-1)
	}
}
