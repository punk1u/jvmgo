package base

import "jvmgo/ch06/rtda"

// 表示指令处理的接口
type Instruction interface {
	// 从字节码中提取操作数
	FetchOperands(reader *BytecodeReader)
	// 执行指令逻辑
	Execute(frame *rtda.Frame)
}

// 表示没有操作数的指令
type NoOperandsInstruction struct {
	// empty
}

// 因为NoOperandsInstruction是没有操作数的指令，没有任何字段，所以FetchOperands提取操作数的方法也为空
func (self *NoOperandsInstruction) FetchOperands(reader *BytecodeReader) {
	// nothing to do
}

// 表示跳转指令，Offset字段存放跳转偏移量
type BranchInstruction struct {
	Offset int
}

// 从字节码中读取一个uint16整数，转成int后赋给Offset字段
func (self *BranchInstruction) FetchOperands(reader *BytecodeReader) {
	self.Offset = int(reader.ReadInt16())
}

type Index8Instruction struct {
	// 存储和加载类指令需要根据索引存取局部变量表，索引由单字节操作数给出。
	// Index字段表示局部变量表索引。
	Index uint
}

func (self *Index8Instruction) FetchOperands(reader *BytecodeReader) {
	self.Index = uint(reader.ReadUint8())
}

// 有一些指令需要访问运行时常量池，常量池索引由两字节操作数给出。
type Index16Instruction struct {
	Index uint
}

func (self *Index16Instruction) FetchOperands(reader *BytecodeReader) {
	self.Index = uint(reader.ReadUint16())
}
