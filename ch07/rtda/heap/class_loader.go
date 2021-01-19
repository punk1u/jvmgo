package heap

import "fmt"
import "jvmgo/ch07/classfile"
import "jvmgo/ch07/classpath"

// 表示类加载器


/*
class names:
    - primitive types: boolean, byte, int ...
    - primitive arrays: [Z, [B, [I ...
    - non-array classes: java/lang/Object ...
    - array classes: [Ljava/lang/Object; ...
*/
type ClassLoader struct {
	// ClassLoader依赖Classpath来搜索和读取class文件，cp字段保存Classpath指针
	cp          *classpath.Classpath
	// 是否打印出类加载信息
	verboseFlag bool
	// 记录已经加载的类数据，key是类的完全限定名
	classMap    map[string]*Class // loaded classes
}

// 创建ClassLoader实例
func NewClassLoader(cp *classpath.Classpath, verboseFlag bool) *ClassLoader {
	return &ClassLoader{
		cp:          cp,
		verboseFlag: verboseFlag,
		classMap:    make(map[string]*Class),
	}
}

// 把类数据加载到方法区
func (self *ClassLoader) LoadClass(name string) *Class {
	// 先查找classMap，看类是否已经被加载。如果是，直接返回类数据，否则调用loadNonArrayClass（）方法加载类
	if class, ok := self.classMap[name]; ok {
		// already loaded
		return class
	}

	return self.loadNonArrayClass(name)
}

/**
数组类和普通类有很大的不同，它的数据并不是来自class文件，
而是由Java虚拟机在运行期间生成。
*/
/**
类的加载大致可以分为三个步骤：首先找到class文件并把数据读取到内存；
然后解析class文件，生成虚拟机可以使用的类数据，并放入方法区；最后进行链接。
**/
func (self *ClassLoader) loadNonArrayClass(name string) *Class {
	data, entry := self.readClass(name)
	class := self.defineClass(data)
	link(class)

	if self.verboseFlag {
		fmt.Printf("[Loaded %s from %s]\n", name, entry)
	}

	return class
}

// 读取class文件中的数据
// 为了打印类加载信息，把最终加载class文件的类路径项也返回给了调用者。
func (self *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
	data, entry, err := self.cp.ReadClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException: " + name)
	}
	return data, entry
}

// 首先调用parseClass（）函数把class文件数据转换成Class结构体。
//Class结构体的superClass和interfaces字段存放超类名和直接接口表，这些类名其实都是符号引用。
// 调用resolveSuperClass（）和resolveInterfaces（）函数解析这些类符号引用
// jvms 5.3.5
func (self *ClassLoader) defineClass(data []byte) *Class {
	class := parseClass(data)
	class.loader = self
	resolveSuperClass(class)
	resolveInterfaces(class)
	self.classMap[class.name] = class
	return class
}

func parseClass(data []byte) *Class {
	cf, err := classfile.Parse(data)
	if err != nil {
		//panic("java.lang.ClassFormatError")
		panic(err)
	}
	return newClass(cf)
}

// 处理父类
// jvms 5.4.3.1
func resolveSuperClass(class *Class) {
	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}

// 处理接口
func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}

// 链接类
func link(class *Class) {
	// 验证类
	verify(class)
	// 准备类
	prepare(class)
}

func verify(class *Class) {
	// todo
}

// jvms 5.4.2
func prepare(class *Class) {
	calcInstanceFieldSlotIds(class)
	calcStaticFieldSlotIds(class)
	allocAndInitStaticVars(class)
}

// 计算实例字段的个数，同时给它们编号
func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	if class.superClass != nil {
		slotId = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.instanceSlotCount = slotId
}

// 计算静态字段的个数，同时给它们编号
func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}

// 给类变量分配空间，然后给它们赋予初始值
func allocAndInitStaticVars(class *Class) {
	class.staticVars = newSlots(class.staticSlotCount)
	for _, field := range class.fields {
		if field.IsStatic() && field.IsFinal() {
			initStaticFinalVar(class, field)
		}
	}
}

// 如果静态变量属于基本类型或String类型，有final修饰符，
// 且它的值在编译期已知，则该值存储在class文件常量池中。

// 从常量池中加载常量值，然后给静态变量赋值
func initStaticFinalVar(class *Class, field *Field) {
	vars := class.staticVars
	cp := class.constantPool
	cpIndex := field.ConstValueIndex()
	slotId := field.SlotId()

	if cpIndex > 0 {
		switch field.Descriptor() {
		case "Z", "B", "C", "S", "I":
			val := cp.GetConstant(cpIndex).(int32)
			vars.SetInt(slotId, val)
		case "J":
			val := cp.GetConstant(cpIndex).(int64)
			vars.SetLong(slotId, val)
		case "F":
			val := cp.GetConstant(cpIndex).(float32)
			vars.SetFloat(slotId, val)
		case "D":
			val := cp.GetConstant(cpIndex).(float64)
			vars.SetDouble(slotId, val)
		case "Ljava/lang/String;":
			panic("todo")
		}
	}
}
