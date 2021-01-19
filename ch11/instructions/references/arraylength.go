package references

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"

// arraylength指令用于获取数组长度。
// arraylength指令只需要一个操作数，即从操作数栈顶弹出的数组引用。
// Get length of array
type ARRAY_LENGTH struct{ base.NoOperandsInstruction }

func (self *ARRAY_LENGTH) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	arrRef := stack.PopRef()
	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arrLen := arrRef.ArrayLength()
	stack.PushInt(arrLen)
}
