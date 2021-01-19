package heap

import "jvmgo/ch10/classfile"

// 表示接口方法符号引用
type InterfaceMethodRef struct {
	MemberRef
	// 缓存解析后的方法指针
	method *Method
}

func newInterfaceMethodRef(cp *ConstantPool, refInfo *classfile.ConstantInterfaceMethodrefInfo) *InterfaceMethodRef {
	ref := &InterfaceMethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

func (self *InterfaceMethodRef) ResolvedInterfaceMethod() *Method {
	if self.method == nil {
		self.resolveInterfaceMethodRef()
	}
	return self.method
}

// 解析接口方法符号引用
// jvms8 5.4.3.4
func (self *InterfaceMethodRef) resolveInterfaceMethodRef() {
	d := self.cp.class
	c := self.ResolvedClass()
	if !c.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	method := lookupInterfaceMethod(c, self.name, self.descriptor)
	if method == nil {
		panic("java.lang.NoSuchMethodError")
	}
	if !method.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	self.method = method
}

// todo
func lookupInterfaceMethod(iface *Class, name, descriptor string) *Method {
	for _, method := range iface.methods {
		if method.name == name && method.descriptor == descriptor {
			return method
		}
	}

	return lookupMethodInInterfaces(iface.interfaces, name, descriptor)
}
