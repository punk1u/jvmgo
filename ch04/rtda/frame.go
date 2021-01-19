package rtda

// stack frame
type Frame struct {
	// stack is implemented as linked list 
	lower        *Frame 
	// save pointer of local vars table
	localVars    LocalVars
	// save Operand Stack 
	operandStack *OperandStack
	// todo
}

func NewFrame(maxLocals, maxStack uint) *Frame {
	return &Frame{
		localVars:    newLocalVars(maxLocals),
		operandStack: newOperandStack(maxStack),
	}
}

// getters
func (self *Frame) LocalVars() LocalVars {
	return self.localVars
}
func (self *Frame) OperandStack() *OperandStack {
	return self.operandStack
}
