package lang

import "unsafe"
import "jvmgo/ch09/native"
import "jvmgo/ch09/rtda"

const jlObject = "java/lang/Object"

func init() {
	native.Register(jlObject, "getClass", "()Ljava/lang/Class;", getClass)
	native.Register(jlObject, "hashCode", "()I", hashCode)
	native.Register(jlObject, "clone", "()Ljava/lang/Object;", clone)
}

/**
首先，从局部变量表中拿到this引用。GetThis（）方法其实就是调用
GetRef（0），不过为了提高代码的可读性，给LocalVars结构体添加了
这个方法。有了this引用后，通过Class（）方法拿到它的Class结构体
指针，进而又通过JClass（）方法拿到它的类对象。最后，把类对象推
入操作数栈顶。
**/
// public final native Class<?> getClass();
// ()Ljava/lang/Class;
func getClass(frame *rtda.Frame) {
	this := frame.LocalVars().GetThis()
	class := this.Class().JClass()
	frame.OperandStack().PushRef(class)
}


// public native int hashCode();
// ()I
func hashCode(frame *rtda.Frame) {
	// 把对象引用（Object结构体指针）转换成uintptr类型，然后强制转换成int32推入操作数栈顶
	this := frame.LocalVars().GetThis()
	hash := int32(uintptr(unsafe.Pointer(this)))
	frame.OperandStack().PushInt(hash)
}

// protected native Object clone() throws CloneNotSupportedException;
// ()Ljava/lang/Object;
func clone(frame *rtda.Frame) {
	this := frame.LocalVars().GetThis()

	/**
	如果类没有实现Cloneable接口，则抛出
	CloneNotSupportedException异常，否则调用Object结构体的
	Clone（）方法克隆对象，然后把对象副本引用推入操作数栈顶。
	**/
	cloneable := this.Class().Loader().LoadClass("java/lang/Cloneable")
	if !this.Class().IsImplements(cloneable) {
		panic("java.lang.CloneNotSupportedException")
	}

	frame.OperandStack().PushRef(this.Clone())
}
