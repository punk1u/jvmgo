package heap

// 表示数组类相关处理方法


// 判断是否是数组对象
func (self *Class) IsArray() bool {
	return self.name[0] == '['
}

// 返回数组类的元素类型
// 先根据数组类名推测出数组元素类名，
// 然后用类加载器加载元素类即可。
func (self *Class) ComponentClass() *Class {
	componentClassName := getComponentClassName(self.name)
	return self.loader.LoadClass(componentClassName)
}

// 根据数组类型创建新的数组对象
func (self *Class) NewArray(count uint) *Object {
	// 如果类并不是数组类，就调用panic（）函数终止程序执行
	if !self.IsArray() {
		panic("Not array class: " + self.name)
	}
	switch self.Name() {
	case "[Z":
		return &Object{self, make([]int8, count)}
	case "[B":
		return &Object{self, make([]int8, count)}
	case "[C":
		return &Object{self, make([]uint16, count)}
	case "[S":
		return &Object{self, make([]int16, count)}
	case "[I":
		return &Object{self, make([]int32, count)}
	case "[J":
		return &Object{self, make([]int64, count)}
	case "[F":
		return &Object{self, make([]float32, count)}
	case "[D":
		return &Object{self, make([]float64, count)}
	default:
		return &Object{self, make([]*Object, count)}
	}
}
