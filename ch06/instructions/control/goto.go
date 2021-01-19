package control

import "jvmgo/ch06/instructions/base"
import "jvmgo/ch06/rtda"

// 控制指令共有11条。jsr和ret指令在Java 6之前用于实现finally子句，
// 从Java 6开始，Oracle的Java编译器已经不再使用这两条指令了，
// 这里不讨论这两条指令。return系列指令有6条，用于从方法调用中返回。
// 剩下的3条指令为：goto、tableswitch和lookupswitch。

// goto指令进行无条件跳转

// Branch always
type GOTO struct{ base.BranchInstruction }

func (self *GOTO) Execute(frame *rtda.Frame) {
	base.Branch(frame, self.Offset)
}
