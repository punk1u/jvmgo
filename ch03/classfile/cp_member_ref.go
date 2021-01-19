package classfile


/**
CONSTANT_Fieldref_info表示字段符号引用，CONSTANT_Methodref_info表示普通（非接口）方法符号引用，CONSTANT_InterfaceMethodref_info表接口方法符号引用。
这三种常量结构一模一样，因此统一使用本结构体表示这三种常量
**/
/*
CONSTANT_Fieldref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
CONSTANT_Methodref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
CONSTANT_InterfaceMethodref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
*/

// Go语言并没有“继承”这个概念，但是可以通过结构体嵌套来模拟，即组合模式。
// 下面三个结构体分别表示CONSTANT_Fieldref_info、CONSTANT_Methodref_info、CONSTANT_InterfaceMethodref_info
type ConstantFieldrefInfo struct{ ConstantMemberrefInfo }
type ConstantMethodrefInfo struct{ ConstantMemberrefInfo }
type ConstantInterfaceMethodrefInfo struct{ ConstantMemberrefInfo }

type ConstantMemberrefInfo struct {
	cp               ConstantPool
    // class_index和name_and_type_index都是常量池索引，分别指向CONSTANT_Class_info和CONSTANT_NameAndType_info常量。
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (self *ConstantMemberrefInfo) readInfo(reader *ClassReader) {
	self.classIndex = reader.readUint16()
	self.nameAndTypeIndex = reader.readUint16()
}

func (self *ConstantMemberrefInfo) ClassName() string {
	return self.cp.getClassName(self.classIndex)
}
func (self *ConstantMemberrefInfo) NameAndDescriptor() (string, string) {
	return self.cp.getNameAndType(self.nameAndTypeIndex)
}
