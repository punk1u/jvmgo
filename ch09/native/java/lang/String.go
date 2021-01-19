package lang

import "jvmgo/ch09/native"
import "jvmgo/ch09/rtda"
import "jvmgo/ch09/rtda/heap"

const jlString = "java/lang/String"

func init() {
	native.Register(jlString, "intern", "()Ljava/lang/String;", intern)
}

// 实现String类的intern（）方法，让Java类库也可以使用它
// public native String intern();
// ()Ljava/lang/String;
func intern(frame *rtda.Frame) {
	this := frame.LocalVars().GetThis()
	interned := heap.InternString(this)
	frame.OperandStack().PushRef(interned)
}
