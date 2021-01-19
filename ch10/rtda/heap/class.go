package heap

import "strings"
import "jvmgo/ch10/classfile"

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
	// 存放Class对应的的文件名
	sourceFile        string
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
	// 表示类的初始化方法是否已经开始执行
	initStarted       bool
	// 表示与Class结构体实例关联的一个类对象(class object)
	jClass            *Object
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
	class.sourceFile = getSourceFile(cf)
	return class
}

// 获取Class对应的文件名
func getSourceFile(cf *classfile.ClassFile) string {
	if sfAttr := cf.SourceFileAttribute(); sfAttr != nil {
		return sfAttr.FileName()
	}
	// 并不是每个class文件中都有源文件信息，因编译器选项而定
	return "Unknown" // todo
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
func (self *Class) Name() string {
	return self.name
}
func (self *Class) ConstantPool() *ConstantPool {
	return self.constantPool
}
func (self *Class) Fields() []*Field {
	return self.fields
}
func (self *Class) Methods() []*Method {
	return self.methods
}
func (self *Class) SourceFile() string {
	return self.sourceFile
}
func (self *Class) Loader() *ClassLoader {
	return self.loader
}
func (self *Class) SuperClass() *Class {
	return self.superClass
}
func (self *Class) StaticVars() Slots {
	return self.staticVars
}
func (self *Class) InitStarted() bool {
	return self.initStarted
}
func (self *Class) JClass() *Object {
	return self.jClass
}

func (self *Class) StartInit() {
	self.initStarted = true
}

// 判断一个Class是否有权限访问另一个Class
// 如果类D想访问类C，需要满足两个条件之一：
// C是public，或者C和D在同一个运行时包内
// jvms 5.4.4
func (self *Class) isAccessibleTo(other *Class) bool {
	return self.IsPublic() ||
		self.GetPackageName() == other.GetPackageName()
}

// 获取包名
// 如果类定义在默认包中，它的包名是空字符串
func (self *Class) GetPackageName() string {
	if i := strings.LastIndex(self.name, "/"); i >= 0 {
		return self.name[:i]
	}
	return ""
}

// 获取类的初始化方法
func (self *Class) GetMainMethod() *Method {
	return self.getMethod("main", "([Ljava/lang/String;)V", true)
}
func (self *Class) GetClinitMethod() *Method {
	return self.getMethod("<clinit>", "()V", true)
}

func (self *Class) getMethod(name, descriptor string, isStatic bool) *Method {
	for c := self; c != nil; c = c.superClass {
		for _, method := range c.methods {
			if method.IsStatic() == isStatic &&
				method.name == name &&
				method.descriptor == descriptor {

				return method
			}
		}
	}
	return nil
}


// 根据字段名和描述符查找字段
func (self *Class) getField(name, descriptor string, isStatic bool) *Field {
	for c := self; c != nil; c = c.superClass {
		for _, field := range c.fields {
			if field.IsStatic() == isStatic &&
				field.name == name &&
				field.descriptor == descriptor {

				return field
			}
		}
	}
	return nil
}

func (self *Class) isJlObject() bool {
	return self.name == "java/lang/Object"
}
func (self *Class) isJlCloneable() bool {
	return self.name == "java/lang/Cloneable"
}
func (self *Class) isJioSerializable() bool {
	return self.name == "java/io/Serializable"
}

func (self *Class) NewObject() *Object {
	return newObject(self)
}

// 返回与类对应的数组类
func (self *Class) ArrayClass() *Class {
	// 根据类名得到数组类名
	arrayClassName := getArrayClassName(self.name)
	// 调用类加载器加载数组类
	return self.loader.LoadClass(arrayClassName)
}

func (self *Class) JavaName() string {
	return strings.Replace(self.name, "/", ".", -1)
}

func (self *Class) IsPrimitive() bool {
	_, ok := primitiveTypes[self.name]
	return ok
}

func (self *Class) GetInstanceMethod(name, descriptor string) *Method {
	return self.getMethod(name, descriptor, false)
}

func (self *Class) GetRefVar(fieldName, fieldDescriptor string) *Object {
	field := self.getField(fieldName, fieldDescriptor, true)
	return self.staticVars.GetRef(field.slotId)
}
func (self *Class) SetRefVar(fieldName, fieldDescriptor string, ref *Object) {
	field := self.getField(fieldName, fieldDescriptor, true)
	self.staticVars.SetRef(field.slotId, ref)
}
