package references

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

/**
putstatic指令给类的某个静态变量赋值，它需要两个操作数。
第一个操作数是uint16索引，来自字节码。通过这个索引可以从当
前类的运行时常量池中找到一个字段符号引用，解析这个符号引用
就可以知道要给类的哪个静态变量赋值。第二个操作数是要赋给静
态变量的值，从操作数栈中弹出。
**/

// Set static field in class
type PUT_STATIC struct{ base.Index16Instruction }

func (self *PUT_STATIC) Execute(frame *rtda.Frame) {
	currentMethod := frame.Method()
	currentClass := currentMethod.Class()
	cp := currentClass.ConstantPool()
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

	/**
	如果解析后的字段是实例字段而非静态字段，则抛出
	IncompatibleClassChangeError异常。如果是final字段，则实际操作的
	是静态常量，只能在类初始化方法中给它赋值。否则，会抛出
	IllegalAccessError异常。
	**/
	if !field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	if field.IsFinal() {
		if currentClass != class || currentMethod.Name() != "<clinit>" {
			panic("java.lang.IllegalAccessError")
		}
	}

	descriptor := field.Descriptor()
	slotId := field.SlotId()
	slots := class.StaticVars()
	stack := frame.OperandStack()

	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		slots.SetInt(slotId, stack.PopInt())
	case 'F':
		slots.SetFloat(slotId, stack.PopFloat())
	case 'J':
		slots.SetLong(slotId, stack.PopLong())
	case 'D':
		slots.SetDouble(slotId, stack.PopDouble())
	case 'L', '[':
		slots.SetRef(slotId, stack.PopRef())
	default:
		// todo
	}
}
