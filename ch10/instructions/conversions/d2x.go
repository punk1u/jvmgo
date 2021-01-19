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

// Convert double to float
type D2F struct{ base.NoOperandsInstruction }

func (self *D2F) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopDouble()
	f := float32(d)
	stack.PushFloat(f)
}

// Convert double to int
type D2I struct{ base.NoOperandsInstruction }

func (self *D2I) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopDouble()
	i := int32(d)
	stack.PushInt(i)
}

// Convert double to long
type D2L struct{ base.NoOperandsInstruction }

func (self *D2L) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopDouble()
	l := int64(d)
	stack.PushLong(l)
}
