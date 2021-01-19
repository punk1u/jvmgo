package reserved

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/native"
import _ "jvmgo/ch11/native/java/io"
import _ "jvmgo/ch11/native/java/lang"
import _ "jvmgo/ch11/native/java/security"
import _ "jvmgo/ch11/native/java/util/concurrent/atomic"
import _ "jvmgo/ch11/native/sun/io"
import _ "jvmgo/ch11/native/sun/misc"
import _ "jvmgo/ch11/native/sun/reflect"

// 调用本地方法指令
// Invoke native method
type INVOKE_NATIVE struct{ base.NoOperandsInstruction }

/**
根据类名、方法名和方法描述符从本地方法注册表中查找本
地方法实现，如果找不到，则抛出UnsatisfiedLinkError异常，否则直
接调用本地方法。
**/
func (self *INVOKE_NATIVE) Execute(frame *rtda.Frame) {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	methodDescriptor := method.Descriptor()

	nativeMethod := native.FindNativeMethod(className, methodName, methodDescriptor)
	if nativeMethod == nil {
		methodInfo := className + "." + methodName + methodDescriptor
		panic("java.lang.UnsatisfiedLinkError: " + methodInfo)
	}

	nativeMethod(frame)
}
