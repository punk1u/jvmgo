package constants

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"

// bipush指令从操作数中获取一个byte型整数，扩展成int型，然后推入栈顶。

// Push byte
type BIPUSH struct {
	val int8
}

func (self *BIPUSH) FetchOperands(reader *base.BytecodeReader) {
	self.val = reader.ReadInt8()
}
func (self *BIPUSH) Execute(frame *rtda.Frame) {
	i := int32(self.val)
	frame.OperandStack().PushInt(i)
}

// sipush指令从操作数中获取一个short型整数，扩展成int型，然后推入栈顶。

// Push short
type SIPUSH struct {
	val int16
}

func (self *SIPUSH) FetchOperands(reader *base.BytecodeReader) {
	self.val = reader.ReadInt16()
}
func (self *SIPUSH) Execute(frame *rtda.Frame) {
	i := int32(self.val)
	frame.OperandStack().PushInt(i)
}
