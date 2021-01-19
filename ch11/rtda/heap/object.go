package heap

// 表示Java中的Object对象
type Object struct {
	// 存放对象的Class指针
	class *Class
	// 存放实例变量(包括普通对象和数组)
	// Go语言的interface{}类型很像C语言中的void*，该类型的变量可以容纳任何类型的值。
	// 对于普通对象来说，data字段中存放的仍然还是Slots变量。但是对于数组，可以在其中放各种类型的数组
	data  interface{} // Slots for Object, []int32 for int[] ...
	// 记录Object结构体实例的额外信息
	extra interface{}
}

// 创建普通非数组对象
// create normal (non-array) object
func newObject(class *Class) *Object {
	return &Object{
		class: class,
		data:  newSlots(class.instanceSlotCount),
	}
}

// getters & setters
func (self *Object) Class() *Class {
	return self.class
}
func (self *Object) Data() interface{} {
	return self.data
}
func (self *Object) Fields() Slots {
	return self.data.(Slots)
}
func (self *Object) Extra() interface{} {
	return self.extra
}
func (self *Object) SetExtra(extra interface{}) {
	self.extra = extra
}

func (self *Object) IsInstanceOf(class *Class) bool {
	return class.IsAssignableFrom(self.class)
}

// reflection
func (self *Object) GetRefVar(name, descriptor string) *Object {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	return slots.GetRef(field.slotId)
}

// 直接给对象的引用类型实例变量赋值
func (self *Object) SetRefVar(name, descriptor string, ref *Object) {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	slots.SetRef(field.slotId, ref)
}
func (self *Object) SetIntVar(name, descriptor string, val int32) {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	slots.SetInt(field.slotId, val)
}
func (self *Object) GetIntVar(name, descriptor string) int32 {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	return slots.GetInt(field.slotId)
}
