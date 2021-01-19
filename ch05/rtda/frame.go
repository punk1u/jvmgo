package rtda

// stack frame
type Frame struct {
	// stack is implemented as linked list 
	lower        *Frame 
	// save pointer of local vars table
	localVars    LocalVars
	// save Operand Stack 
	operandStack *OperandStack
	thread       *Thread
	nextPC       int // the next instruction after the call
}

func newFrame(thread *Thread, maxLocals, maxStack uint) *Frame {
	return &Frame{
		thread:       thread,
		localVars:    newLocalVars(maxLocals),
		operandStack: newOperandStack(maxStack),
	}
}

// getters & setters
func (self *Frame) LocalVars() LocalVars {
	return self.localVars
}
func (self *Frame) OperandStack() *OperandStack {
	return self.operandStack
}
func (self *Frame) Thread() *Thread {
	return self.thread
}
func (self *Frame) NextPC() int {
	return self.nextPC
}
func (self *Frame) SetNextPC(nextPC int) {
	self.nextPC = nextPC
}
