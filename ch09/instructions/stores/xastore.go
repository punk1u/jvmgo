package stores

import "jvmgo/ch09/instructions/base"
import "jvmgo/ch09/rtda"
import "jvmgo/ch09/rtda/heap"

// <t>astore系列指令按索引给数组元素赋值。
// Store into reference array
type AASTORE struct{ base.NoOperandsInstruction }

/**
iastore指令的三个操作数分别是：要赋给数组元素的值、数组
索引、数组引用，依次从操作数栈中弹出。如果数组引用是null，则
抛出NullPointerException。如果数组索引小于0或者大于等于数组
长度，则抛出ArrayIndexOutOfBoundsException异常。这两个检查和
<t>aload系列指令一样。如果一切正常，则按索引给数组元素赋值。
**/
func (self *AASTORE) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	ref := stack.PopRef()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	refs := arrRef.Refs()
	checkIndex(len(refs), index)
	refs[index] = ref
}

// Store into byte or boolean array
type BASTORE struct{ base.NoOperandsInstruction }

func (self *BASTORE) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	bytes := arrRef.Bytes()
	checkIndex(len(bytes), index)
	bytes[index] = int8(val)
}

// Store into char array
type CASTORE struct{ base.NoOperandsInstruction }

func (self *CASTORE) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	chars := arrRef.Chars()
	checkIndex(len(chars), index)
	chars[index] = uint16(val)
}

// Store into double array
type DASTORE struct{ base.NoOperandsInstruction }

func (self *DASTORE) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	val := stack.PopDouble()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	doubles := arrRef.Doubles()
	checkIndex(len(doubles), index)
	doubles[index] = float64(val)
}

// Store into float array
type FASTORE struct{ base.NoOperandsInstruction }

func (self *FASTORE) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	val := stack.PopFloat()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	floats := arrRef.Floats()
	checkIndex(len(floats), index)
	floats[index] = float32(val)
}

// Store into int array
type IASTORE struct{ base.NoOperandsInstruction }

func (self *IASTORE) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	ints := arrRef.Ints()
	checkIndex(len(ints), index)
	ints[index] = int32(val)
}

// Store into long array
type LASTORE struct{ base.NoOperandsInstruction }

func (self *LASTORE) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	val := stack.PopLong()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	longs := arrRef.Longs()
	checkIndex(len(longs), index)
	longs[index] = int64(val)
}

// Store into short array
type SASTORE struct{ base.NoOperandsInstruction }

func (self *SASTORE) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	shorts := arrRef.Shorts()
	checkIndex(len(shorts), index)
	shorts[index] = int16(val)
}

func checkNotNil(ref *heap.Object) {
	if ref == nil {
		panic("java.lang.NullPointerException")
	}
}
func checkIndex(arrLen int, index int32) {
	if index < 0 || index >= int32(arrLen) {
		panic("ArrayIndexOutOfBoundsException")
	}
}
