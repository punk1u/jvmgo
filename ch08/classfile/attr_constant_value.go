package classfile

/**
表示ConstantValue定长属性，只会出现在field_info结构中，用于表示常量表达式的值
 下表为字段类型和常量类型的关系：

| 字段类型                    | 常量类型              |
| --------------------------- | --------------------- |
| long                        | CONSTANT_Long_info    |
| float                       | CONSTANT_Float_info   |
| double                      | CONSTANT_Double_info  |
| int,short,char,byte,boolean | CONSTANT_Integer_info |
| String                      | CONSTANT_String_info  |
讨论类和对象时，会介绍如何使用ConstantValue属性。
**/
/*
ConstantValue_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 constantvalue_index;
}
*/
type ConstantValueAttribute struct {
	constantValueIndex uint16
}

func (self *ConstantValueAttribute) readInfo(reader *ClassReader) {
	self.constantValueIndex = reader.readUint16()
}

func (self *ConstantValueAttribute) ConstantValueIndex() uint16 {
	return self.constantValueIndex
}
