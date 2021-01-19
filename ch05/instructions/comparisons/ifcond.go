package comparisons

import "jvmgo/ch05/instructions/base"
import "jvmgo/ch05/rtda"


/**
if<cond>指令把操作数栈顶的int变量弹出，然后跟0进行比较，满足条件则跳转。

假设从栈顶弹出的变量是x，则指令执行跳转操作的条件如下：

1. ifeq：x==0
2. ifne：x！=0
3. iflt：x<0
4. ifle：x<=0
5. ifgt：x>0
6. ifge：x>=0

真正的跳转逻辑在Branch（）函数中。因为这个函数在很多指令中都会用到，所以抽离出去
**/



// Branch if int comparison with zero succeeds
type IFEQ struct{ base.BranchInstruction }

func (self *IFEQ) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val == 0 {
		base.Branch(frame, self.Offset)
	}
}

type IFNE struct{ base.BranchInstruction }

func (self *IFNE) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val != 0 {
		base.Branch(frame, self.Offset)
	}
}

type IFLT struct{ base.BranchInstruction }

func (self *IFLT) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val < 0 {
		base.Branch(frame, self.Offset)
	}
}

type IFLE struct{ base.BranchInstruction }

func (self *IFLE) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val <= 0 {
		base.Branch(frame, self.Offset)
	}
}

type IFGT struct{ base.BranchInstruction }

func (self *IFGT) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val > 0 {
		base.Branch(frame, self.Offset)
	}
}

type IFGE struct{ base.BranchInstruction }

func (self *IFGE) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val >= 0 {
		base.Branch(frame, self.Offset)
	}
}
