package rtda

import "jvmgo/ch06/rtda/heap"

/**
局部变量表是按索引访问的，所以很自然，可以把它想象成一
个数组。根据Java虚拟机规范，这个数组的每个元素至少可以容纳
一个int或引用值，两个连续的元素可以容纳一个long或double值。

此结构体表示可以同时容纳一个int值和一个引用值的对象
**/
type Slot struct {
	num int32
	ref *heap.Object
}
