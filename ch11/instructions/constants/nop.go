package constants

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"

// nop指令是最简单的一条指令，因为它什么也不做。
// Do nothing
type NOP struct{ base.NoOperandsInstruction }

func (self *NOP) Execute(frame *rtda.Frame) {
	// really do nothing
}
