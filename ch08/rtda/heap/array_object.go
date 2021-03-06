package heap

// 针对基本类型Byte数组返回具体的数组数据
func (self *Object) Bytes() []int8 {
	return self.data.([]int8)
}

// 针对基本类型Short数组返回具体的数组数据
func (self *Object) Shorts() []int16 {
	return self.data.([]int16)
}

// 针对基本类型Int数组返回具体的数组数据
func (self *Object) Ints() []int32 {
	return self.data.([]int32)
}

// 针对基本类型Long数组返回具体的数组数据
func (self *Object) Longs() []int64 {
	return self.data.([]int64)
}

// 针对基本类型Char数组返回具体的数组数据
func (self *Object) Chars() []uint16 {
	return self.data.([]uint16)
}

// 针对基本类型Float数组返回具体的数组数据
func (self *Object) Floats() []float32 {
	return self.data.([]float32)
}

// 针对基本类型Double数组返回具体的数组数据
func (self *Object) Doubles() []float64 {
	return self.data.([]float64)
}

// 针对引用类型数组返回具体的数组数据
func (self *Object) Refs() []*Object {
	return self.data.([]*Object)
}

// 获取数组长度
func (self *Object) ArrayLength() int32 {
	switch self.data.(type) {
	case []int8:
		return int32(len(self.data.([]int8)))
	case []int16:
		return int32(len(self.data.([]int16)))
	case []int32:
		return int32(len(self.data.([]int32)))
	case []int64:
		return int32(len(self.data.([]int64)))
	case []uint16:
		return int32(len(self.data.([]uint16)))
	case []float32:
		return int32(len(self.data.([]float32)))
	case []float64:
		return int32(len(self.data.([]float64)))
	case []*Object:
		return int32(len(self.data.([]*Object)))
	default:
		panic("Not array!")
	}
}
