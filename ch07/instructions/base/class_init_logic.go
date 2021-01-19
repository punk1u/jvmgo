package base

import "jvmgo/ch07/rtda"
import "jvmgo/ch07/rtda/heap"

/**
表示类初始化方法
InitClass（）函数先调用StartInit（）方法把类的initStarted状态设
置成true以免进入死循环，然后调用scheduleClinit（）函数准备执行
类的初始化方法
**/
// jvms 5.5
func InitClass(thread *rtda.Thread, class *heap.Class) {
	class.StartInit()
	scheduleClinit(thread, class)
	initSuperClass(thread, class)
}

func scheduleClinit(thread *rtda.Thread, class *heap.Class) {
	clinit := class.GetClinitMethod()
	if clinit != nil {
		// exec <clinit>
		newFrame := thread.NewFrame(clinit)
		thread.PushFrame(newFrame)
	}
}

/**
如果超类的初始化还没有开始，就递归调用InitClass（）函数执
行超类的初始化方法，这样可以保证超类的初始化方法对应的帧在
子类上面，使超类初始化方法先于子类执行。
**/
func initSuperClass(thread *rtda.Thread, class *heap.Class) {
	if !class.IsInterface() {
		superClass := class.SuperClass()
		if superClass != nil && !superClass.InitStarted() {
			InitClass(thread, superClass)
		}
	}
}
