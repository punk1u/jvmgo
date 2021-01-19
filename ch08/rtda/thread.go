package rtda

import "jvmgo/ch08/rtda/heap"

/*
表示线程的结构体 

线程私有的运行时数据区用于辅助执行Java字节码。每个线程都有自己的pc寄存器（Program Counter）和Java虚拟机栈（JVMStack）。
Java虚拟机栈又由栈帧（Stack Frame，后面简称帧）构成，帧中保存方法执行的状态，
包括局部变量表（Local Variable）和操作数栈（Operand Stack）等。
在任一时刻，某一线程肯定是在执行某个方法。这个方法叫作该线程的当前方法；
执行该方法的帧叫作线程的当前帧；
声明该方法的类叫作当前类。如果当前方法是Java方法，则pc寄存器中存放当前正在执行的Java虚拟机指令的地址，
否则，当前方法是本地方法，pc寄存器中的值没有明确定义。
JVM
  Thread
    pc 寄存器
    Stack 栈
      Frame  栈帧
        LocalVars
        OperandStack
*/
type Thread struct {
	pc    int // the address of the instruction currently being executed
	stack *Stack // the pointer of Stack struct（Jvm stack）
	// todo
}

/**
创建线程
**/
func NewThread() *Thread {
	return &Thread{
		// paramete 1024 means max number of stack frame
		// TODO fix cmd util to set this parameter
		stack: newStack(1024),
	}
}

func (self *Thread) PC() int {
	return self.pc
}
func (self *Thread) SetPC(pc int) {
	self.pc = pc
}

// push the stack frame
func (self *Thread) PushFrame(frame *Frame) {
	self.stack.push(frame)
}

// pop the stack frame
func (self *Thread) PopFrame() *Frame {
	return self.stack.pop()
}

// return current stack frame
func (self *Thread) CurrentFrame() *Frame {
	return self.stack.top()
}
func (self *Thread) TopFrame() *Frame {
	return self.stack.top()
}

func (self *Thread) IsStackEmpty() bool {
	return self.stack.isEmpty()
}

// 创建一个新的栈帧
func (self *Thread) NewFrame(method *heap.Method) *Frame {
	return newFrame(self, method)
}
