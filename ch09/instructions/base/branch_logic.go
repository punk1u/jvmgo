package base

import "jvmgo/ch09/rtda"

/**
跳转逻辑，因为这个跳转逻辑在很多指令中都会用到，所以单独定义出来
**/
func Branch(frame *rtda.Frame, offset int) {
	pc := frame.Thread().PC()
	nextPC := pc + offset
	frame.SetNextPC(nextPC)
}
