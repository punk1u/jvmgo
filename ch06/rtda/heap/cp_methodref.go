package heap

import "jvmgo/ch06/classfile"

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

func (self *MethodRef) ResolvedMethod() *Method {
	if self.method == nil {
		self.resolveMethodRef()
	}
	return self.method
}

// jvms8 5.4.3.3
func (self *MethodRef) resolveMethodRef() {
	//class := self.Class()
	// todo
}
