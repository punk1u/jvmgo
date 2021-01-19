package comparisons

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"

// if_acmpeq和if_acmpne指令把栈顶的两个引用弹出，根据引用是否相同进行跳转。

// Branch if reference comparison succeeds
type IF_ACMPEQ struct{ base.BranchInstruction }

func (self *IF_ACMPEQ) Execute(frame *rtda.Frame) {
	if _acmp(frame) {
		base.Branch(frame, self.Offset)
	}
}

type IF_ACMPNE struct{ base.BranchInstruction }

func (self *IF_ACMPNE) Execute(frame *rtda.Frame) {
	if !_acmp(frame) {
		base.Branch(frame, self.Offset)
	}
}

func _acmp(frame *rtda.Frame) bool {
	stack := frame.OperandStack()
	ref2 := stack.PopRef()
	ref1 := stack.PopRef()
	return ref1 == ref2 // todo
}
