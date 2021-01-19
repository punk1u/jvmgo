package native

import "jvmgo/ch11/rtda"

// 本地方法注册表，用来注册和查找本地方法

// 表示本地方法
// 参数是Frame结构体指针，没有
// 返回值。这个frame参数就是本地方法的工作空间
type NativeMethod func(frame *rtda.Frame)

// 表示注册表
// 类名、方法名和方法描述符加在一起才能唯一确定一个方法，
// 所以把它们的组合作为本地方法注册表的键
// 值是具体的本地方法实现
var registry = map[string]NativeMethod{}

/**
jva.lang.Object等类是通过一个叫作registerNatives（）
的本地方法来注册其他本地方法的。将自己注册所有的
本地方法实现。所以像registerNatives（）这样的方法就没有太大的用
处。为了避免重复代码，这里统一处理，如果遇到这样的本地方
法，就返回一个空的实现
**/
func emptyNativeMethod(frame *rtda.Frame) {
	// do nothing
}

// 注册本地方法
func Register(className, methodName, methodDescriptor string, method NativeMethod) {
	key := className + "~" + methodName + "~" + methodDescriptor
	registry[key] = method
}

// 查找本地方法
func FindNativeMethod(className, methodName, methodDescriptor string) NativeMethod {
	key := className + "~" + methodName + "~" + methodDescriptor
	if method, ok := registry[key]; ok {
		return method
	}
	if methodDescriptor == "()V" {
		if methodName == "registerNatives" || methodName == "initIDs" {
			return emptyNativeMethod
		}
	}
	return nil
}
