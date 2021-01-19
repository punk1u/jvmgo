package classfile

import "math"

// 用于统一表示integer,boolean、byte、short和char

// 表示虚拟机规范中的CONSTANT_Integer_info。CONSTANT_Integer_info正好可以容纳一个Java的int型常量，
// 但实际上比int更小的boolean、byte、short和char也放在CONSTANT_Integer_info中
/*
CONSTANT_Integer_info {
    u1 tag;
    u4 bytes;
}
*/
type ConstantIntegerInfo struct {
	val int32
}


// 读取一个uint32数据，然后把它转型为int32类型
func (self *ConstantIntegerInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint32()
	self.val = int32(bytes)
}
func (self *ConstantIntegerInfo) Value() int32 {
	return self.val
}

// 表示虚拟机规范中的CONSTANT_Float_info单精度浮点数，
/*
CONSTANT_Float_info {
    u1 tag;
    u4 bytes;
}
*/
type ConstantFloatInfo struct {
	val float32
}

// 读取一个uint32数据，然后把它转型为float32类型
func (self *ConstantFloatInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint32()
	self.val = math.Float32frombits(bytes)
}
func (self *ConstantFloatInfo) Value() float32 {
	return self.val
}


// 表示虚拟机规范中的CONSTANT_Long_info
/*
CONSTANT_Long_info {
    u1 tag;
    u4 high_bytes;
    u4 low_bytes;
}
*/
type ConstantLongInfo struct {
	val int64
}

// 先读取一个uint64数据，然后把它转型为int64类型
func (self *ConstantLongInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint64()
	self.val = int64(bytes)
}
func (self *ConstantLongInfo) Value() int64 {
	return self.val
}

// 表示虚拟机规范中的CONSTANT_Double_info
/*
CONSTANT_Double_info {
    u1 tag;
    u4 high_bytes;
    u4 low_bytes;
}
*/
type ConstantDoubleInfo struct {
	val float64
}

// 先读取一个uint64数据，然后调用math包的Float64frombits()函数把它转换成float64类型
func (self *ConstantDoubleInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint64()
	self.val = math.Float64frombits(bytes)
}
func (self *ConstantDoubleInfo) Value() float64 {
	return self.val
}
