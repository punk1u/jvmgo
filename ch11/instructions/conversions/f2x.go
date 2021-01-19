package conversions

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"

/**
类型转换指令大致对应Java语言中的基本类型强制转换操作。类型转换指令有共15条，
引用类型转换对应的是checkcast指令。



按照被转换变量的类型，类型转换指令可以分为4种：
1、i2x系列指令把int变量强制转换成其他类型；
2、l2x系列指令把long变量强制转换成其他类型；
3、f2x系列指令把float变量强制转换成其他类型；
4、d2x系列指令把double变量强制转换成其他类型。
**/

// Convert float to double
type F2D struct{ base.NoOperandsInstruction }

func (self *F2D) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	f := stack.PopFloat()
	d := float64(f)
	stack.PushDouble(d)
}

// Convert float to int
type F2I struct{ base.NoOperandsInstruction }

func (self *F2I) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	f := stack.PopFloat()
	i := int32(f)
	stack.PushInt(i)
}

// Convert float to long
type F2L struct{ base.NoOperandsInstruction }

func (self *F2L) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	f := stack.PopFloat()
	l := int64(f)
	stack.PushLong(l)
}
