package references

import "jvmgo/ch10/instructions/base"
import "jvmgo/ch10/rtda"
import "jvmgo/ch10/rtda/heap"

// getstatic指令和putstatic正好相反，它取出类的某个静态变量值，然后推入栈顶。

// Get static field from class
type GET_STATIC struct{ base.Index16Instruction }

/**
getstatic指令只需要一个操作数：uint16常量池索引
如果解析后的字段不是静态字段，也要抛出
IncompatibleClassChangeError异常。如果声明字段的类还没有初始
化好，也需要先初始化。getstatic只是读取静态变量的值，自然也就
不用管它是否是final了。
**/
func (self *GET_STATIC) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)
	field := fieldRef.ResolvedField()
	class := field.Class()
	/**
	先判
	断类的初始化是否已经开始，如果还没有，则需要调用类的初始化
	方法，并终止指令执行。但是由于此时指令已经执行到了一半，也
	就是说当前帧的nextPC字段已经指向下一条指令，所以需要修改
	nextPC，让它重新指向当前指令。
	**/
	if !class.InitStarted() {
		frame.RevertNextPC()
		base.InitClass(frame.Thread(), class)
		return
	}

	if !field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	descriptor := field.Descriptor()
	slotId := field.SlotId()
	slots := class.StaticVars()
	stack := frame.OperandStack()

	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		stack.PushInt(slots.GetInt(slotId))
	case 'F':
		stack.PushFloat(slots.GetFloat(slotId))
	case 'J':
		stack.PushLong(slots.GetLong(slotId))
	case 'D':
		stack.PushDouble(slots.GetDouble(slotId))
	case 'L', '[':
		stack.PushRef(slots.GetRef(slotId))
	default:
		// todo
	}
}
