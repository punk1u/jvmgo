package heap

import "jvmgo/ch10/classfile"

// 表示方法符号引用
type MethodRef struct {
	MemberRef
	// 缓存解析后的方法指针
	method *Method
}

// 创建MethodRef实例
func newMethodRef(cp *ConstantPool, refInfo *classfile.ConstantMethodrefInfo) *MethodRef {
	ref := &MethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

// 解析符号引用
func (self *MethodRef) ResolvedMethod() *Method {
	if self.method == nil {
		self.resolveMethodRef()
	}
	return self.method
}

/**
如果类D想通过方法符号引用访问类C的某个方法，先要解析
符号引用得到类C。如果C是接口，则抛出
IncompatibleClassChangeError异常，否则根据方法名和描述符查找
方法。如果找不到对应的方法，则抛出NoSuchMethodError异常，否
则检查类D是否有权限访问该方法。如果没有，则抛出
IllegalAccessError异常。
**/
// jvms8 5.4.3.3
func (self *MethodRef) resolveMethodRef() {
	d := self.cp.class
	c := self.ResolvedClass()
	if c.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	method := lookupMethod(c, self.name, self.descriptor)
	if method == nil {
		panic("java.lang.NoSuchMethodError")
	}
	if !method.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	self.method = method
}

// 先从C的继承层次中找，如果找不到，就去C的接口中找。
func lookupMethod(class *Class, name, descriptor string) *Method {
	method := LookupMethodInClass(class, name, descriptor)
	if method == nil {
		method = lookupMethodInInterfaces(class.interfaces, name, descriptor)
	}
	return method
}
