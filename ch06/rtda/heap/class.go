package heap

import "strings"
import "jvmgo/ch06/classfile"

// name, superClassName and interfaceNames are all binary names(jvms8-4.2.1)
type Class struct {
	// 类的访问标志，总共16比特。字段和方法也有访问标志，但具体标志位的含义可能有所不同
	accessFlags       uint16
	// 类名(完全限定名)
	name              string // thisClassName
	// 超类名(完全限定名)
	superClassName    string
	// 接口名(完全限定名)
	interfaceNames    []string
	// 存放运行时常量池指针
	constantPool      *ConstantPool
	// 存放字段表
	fields            []*Field
	// 存放方法表
	methods           []*Method
	// loader字段存放类加载器指针
	loader            *ClassLoader
	// 类的超类
	superClass        *Class
	// 接口指针
	interfaces        []*Class
	// 实例变量占据的空间大小
	instanceSlotCount uint
	// 类变量占据的空间大小
	staticSlotCount   uint
	// 存放静态变量
	staticVars        Slots
}

// 把ClassFile结构体转换成Class结构体
func newClass(cf *classfile.ClassFile) *Class {
	class := &Class{}
	class.accessFlags = cf.AccessFlags()
	class.name = cf.ClassName()
	class.superClassName = cf.SuperClassName()
	class.interfaceNames = cf.InterfaceNames()
	class.constantPool = newConstantPool(class, cf.ConstantPool())
	class.fields = newFields(class, cf.Fields())
	class.methods = newMethods(class, cf.Methods())
	return class
}

// 判断是否是公共的
func (self *Class) IsPublic() bool {
	return 0 != self.accessFlags&ACC_PUBLIC
}
func (self *Class) IsFinal() bool {
	return 0 != self.accessFlags&ACC_FINAL
}
func (self *Class) IsSuper() bool {
	return 0 != self.accessFlags&ACC_SUPER
}
func (self *Class) IsInterface() bool {
	return 0 != self.accessFlags&ACC_INTERFACE
}
func (self *Class) IsAbstract() bool {
	return 0 != self.accessFlags&ACC_ABSTRACT
}
func (self *Class) IsSynthetic() bool {
	return 0 != self.accessFlags&ACC_SYNTHETIC
}
func (self *Class) IsAnnotation() bool {
	return 0 != self.accessFlags&ACC_ANNOTATION
}
func (self *Class) IsEnum() bool {
	return 0 != self.accessFlags&ACC_ENUM
}

// getters
func (self *Class) ConstantPool() *ConstantPool {
	return self.constantPool
}
func (self *Class) StaticVars() Slots {
	return self.staticVars
}

// 判断一个Class是否有权限访问另一个Class
// 如果类D想访问类C，需要满足两个条件之一：
// C是public，或者C和D在同一个运行时包内
// jvms 5.4.4
func (self *Class) isAccessibleTo(other *Class) bool {
	return self.IsPublic() ||
		self.getPackageName() == other.getPackageName()
}

// 获取包名
// 如果类定义在默认包中，它的包名是空字符串
func (self *Class) getPackageName() string {
	if i := strings.LastIndex(self.name, "/"); i >= 0 {
		return self.name[:i]
	}
	return ""
}

func (self *Class) GetMainMethod() *Method {
	return self.getStaticMethod("main", "([Ljava/lang/String;)V")
}

func (self *Class) getStaticMethod(name, descriptor string) *Method {
	for _, method := range self.methods {
		if method.IsStatic() &&
			method.name == name &&
			method.descriptor == descriptor {

			return method
		}
	}
	return nil
}

func (self *Class) NewObject() *Object {
	return newObject(self)
}
