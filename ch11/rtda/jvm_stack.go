package rtda

// jvm stack
// stack is implemented as linked list,so stack can use memory when it needs,and poped stack frame can be collected by Go
type Stack struct {
	// max stack frame number
	maxSize uint
	// current stack frame size
	size    uint
	// save the top stack frame
	_top    *Frame // stack is implemented as linked list
}

func newStack(maxSize uint) *Stack {
	return &Stack{
		maxSize: maxSize,
	}
}

// push new stack frame in the stack
func (self *Stack) push(frame *Frame) {
	if self.size >= self.maxSize {
		panic("java.lang.StackOverflowError")
	}

	if self._top != nil {
		frame.lower = self._top
	}

	self._top = frame
	self.size++
}

func (self *Stack) pop() *Frame {
	if self._top == nil {
		panic("jvm stack is empty!")
	}

	top := self._top
	self._top = top.lower
	top.lower = nil
	self.size--

	return top
}

// return current top stack frame
func (self *Stack) top() *Frame {
	if self._top == nil {
		panic("jvm stack is empty!")
	}

	return self._top
}

// 获取完整的Java虚拟机栈
func (self *Stack) getFrames() []*Frame {
	frames := make([]*Frame, 0, self.size)
	for frame := self._top; frame != nil; frame = frame.lower {
		frames = append(frames, frame)
	}
	return frames
}

func (self *Stack) isEmpty() bool {
	return self._top == nil
}

// 清空线程(Thread)的操作栈
func (self *Stack) clear() {
	for !self.isEmpty() {
		self.pop()
	}
}
