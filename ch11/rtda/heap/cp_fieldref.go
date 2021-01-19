package heap

import "jvmgo/ch11/classfile"

// 表示字段符号引用
type FieldRef struct {
	MemberRef
	// 缓存解析后的字段指针
	field *Field
}

// 创建FieldRef实例
func newFieldRef(cp *ConstantPool, refInfo *classfile.ConstantFieldrefInfo) *FieldRef {
	ref := &FieldRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

func (self *FieldRef) ResolvedField() *Field {
	if self.field == nil {
		self.resolveFieldRef()
	}
	return self.field
}

// 解析字段符号引用
/**
如果类D想通过字段符号引用访问类C的某个字段，首先要解
析符号引用得到类C，然后根据字段名和描述符查找字段。如果字
段查找失败，则虚拟机抛出NoSuchFieldError异常。如果查找成功，
但D没有足够的权限访问该字段，则虚拟机抛出IllegalAccessError异常。
**/
// jvms 5.4.3.2
func (self *FieldRef) resolveFieldRef() {
	d := self.cp.class
	c := self.ResolvedClass()
	field := lookupField(c, self.name, self.descriptor)

	if field == nil {
		panic("java.lang.NoSuchFieldError")
	}
	if !field.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	self.field = field
}

// 在Class中查找对应的字段
/**
首先在C的字段中查找。如果找不到，在C的直接接口递归应
用这个查找过程。如果还找不到的话，在C的超类中递归应用这个
查找过程。如果仍然找不到，则查找失败。
*/
func lookupField(c *Class, name, descriptor string) *Field {
	for _, field := range c.fields {
		if field.name == name && field.descriptor == descriptor {
			return field
		}
	}

	for _, iface := range c.interfaces {
		if field := lookupField(iface, name, descriptor); field != nil {
			return field
		}
	}

	if c.superClass != nil {
		return lookupField(c.superClass, name, descriptor)
	}

	return nil
}
