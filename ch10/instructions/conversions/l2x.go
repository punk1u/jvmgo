package conversions

import "jvmgo/ch10/instructions/base"
import "jvmgo/ch10/rtda"

/**
类型转换指令大致对应Java语言中的基本类型强制转换操作。类型转换指令有共15条，
引用类型转换对应的是checkcast指令。



按照被转换变量的类型，类型转换指令可以分为4种：
1、i2x系列指令把int变量强制转换成其他类型；
2、l2x系列指令把long变量强制转换成其他类型；
3、f2x系列指令把float变量强制转换成其他类型；
4、d2x系列指令把double变量强制转换成其他类型。
**/

// Convert long to double
type L2D struct{ base.NoOperandsInstruction }

func (self *L2D) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	l := stack.PopLong()
	d := float64(l)
	stack.PushDouble(d)
}

// Convert long to float
type L2F struct{ base.NoOperandsInstruction }

func (self *L2F) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	l := stack.PopLong()
	f := float32(l)
	stack.PushFloat(f)
}

// Convert long to int
type L2I struct{ base.NoOperandsInstruction }

func (self *L2I) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	l := stack.PopLong()
	i := int32(l)
	stack.PushInt(i)
}
