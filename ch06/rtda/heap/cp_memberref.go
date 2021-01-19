package heap

import "jvmgo/ch06/classfile"

// 存放字段和方法符号引用共有的信息
/**
站在Java虚拟机的角度，一个类是完全可以有多个同名字段的，只要它们的类型互不相同就可以。
所以字段符号引用要存放字段描述符
**/

type MemberRef struct {
	SymRef
	name       string
	descriptor string
}

// 从class文件内存储的字段或方法常量中提取数据
func (self *MemberRef) copyMemberRefInfo(refInfo *classfile.ConstantMemberrefInfo) {
	self.className = refInfo.ClassName()
	self.name, self.descriptor = refInfo.NameAndDescriptor()
}

func (self *MemberRef) Name() string {
	return self.name
}
func (self *MemberRef) Descriptor() string {
	return self.descriptor
}
