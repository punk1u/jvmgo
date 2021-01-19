package heap

// 符号引用，统一表示类符号引用、字段符号引用、方法符号引用和接口方法符号引用的父类公共数据

// symbolic reference
type SymRef struct {
	// 存放符号引用所在的运行时常量池指针
	cp        *ConstantPool
	// 存放类的完全限定名
	className string
	// 缓存解析后的类结构体指针
	class     *Class
}

func (self *SymRef) ResolvedClass() *Class {
	if self.class == nil {
		self.resolveClassRef()
	}
	return self.class
}

// 解析类符号引用
// 如果类D通过符号引用N引用类C的话，要解析N，先用D的类加载器加载C，然后检查D是否有权限访问C，
// 如果没 有，则抛出IllegalAccessError异常
// jvms8 5.4.3.1
func (self *SymRef) resolveClassRef() {
	d := self.cp.class
	c := d.loader.LoadClass(self.className)
	if !c.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	self.class = c
}
