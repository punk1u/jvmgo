package heap

import "jvmgo/ch10/classfile"

// 表示方法
type Method struct {
	ClassMember
	// 存放操作数栈
	maxStack        uint
	// 局部变量表大小
	maxLocals       uint
	// code字段存放方法字节码
	code            []byte
	// 异常处理表
	exceptionTable  ExceptionTable
	// 表示方法对应的行号
	lineNumberTable *classfile.LineNumberTableAttribute
	// 表示方法参数个数
	argSlotCount    uint
}

// 根据class文件中的方法信息创建Method表
func newMethods(class *Class, cfMethods []*classfile.MemberInfo) []*Method {
	methods := make([]*Method, len(cfMethods))
	for i, cfMethod := range cfMethods {
		methods[i] = newMethod(class, cfMethod)
	}
	return methods
}

func newMethod(class *Class, cfMethod *classfile.MemberInfo) *Method {
	method := &Method{}
	method.class = class
	method.copyMemberInfo(cfMethod)
	method.copyAttributes(cfMethod)
	md := parseMethodDescriptor(method.descriptor)
	method.calcArgSlotCount(md.parameterTypes)
	// 如果是本地方法，则注入字节码和其他信息
	if method.IsNative() {
		method.injectCodeAttribute(md.returnType)
	}
	return method
}

// 从class文件中中拷贝属性
func (self *Method) copyAttributes(cfMethod *classfile.MemberInfo) {
	if codeAttr := cfMethod.CodeAttribute(); codeAttr != nil {
		self.maxStack = codeAttr.MaxStack()
		self.maxLocals = codeAttr.MaxLocals()
		self.code = codeAttr.Code()
		// 从class文件中提取行号表
		self.lineNumberTable = codeAttr.LineNumberTableAttribute()
		// 从Code属性中复制异常处理表
		self.exceptionTable = newExceptionTable(codeAttr.ExceptionTable(),
			self.class.constantPool)
	}
}

func (self *Method) calcArgSlotCount(paramTypes []string) {
	for _, paramType := range paramTypes {
		self.argSlotCount++
		if paramType == "J" || paramType == "D" {
			self.argSlotCount++
		}
	}
	if !self.IsStatic() {
		self.argSlotCount++ // `this` reference
	}
}

// 注入字节码和其他信息
/**
本地方法在class文件中没有Code属性，所以需要给maxStack和
maxLocals字段赋值。本地方法帧的操作数栈至少要能容纳返回值，
为了简化代码，暂时给maxStack字段赋值为4。因为本地方法帧的
局部变量表只用来存放参数值，所以把argSlotCount赋给maxLocals
字段刚好。至于code字段，也就是本地方法的字节码，第一条指令
都是0xFE，第二条指令则根据函数的返回值选择相应的返回指令。
**/
func (self *Method) injectCodeAttribute(returnType string) {
	self.maxStack = 4 // todo
	self.maxLocals = self.argSlotCount
	switch returnType[0] {
	case 'V':
		self.code = []byte{0xfe, 0xb1} // return
	case 'L', '[':
		self.code = []byte{0xfe, 0xb0} // areturn
	case 'D':
		self.code = []byte{0xfe, 0xaf} // dreturn
	case 'F':
		self.code = []byte{0xfe, 0xae} // freturn
	case 'J':
		self.code = []byte{0xfe, 0xad} // lreturn
	default:
		self.code = []byte{0xfe, 0xac} // ireturn
	}
}

func (self *Method) IsSynchronized() bool {
	return 0 != self.accessFlags&ACC_SYNCHRONIZED
}
func (self *Method) IsBridge() bool {
	return 0 != self.accessFlags&ACC_BRIDGE
}
func (self *Method) IsVarargs() bool {
	return 0 != self.accessFlags&ACC_VARARGS
}
func (self *Method) IsNative() bool {
	return 0 != self.accessFlags&ACC_NATIVE
}
func (self *Method) IsAbstract() bool {
	return 0 != self.accessFlags&ACC_ABSTRACT
}
func (self *Method) IsStrict() bool {
	return 0 != self.accessFlags&ACC_STRICT
}

// getters
func (self *Method) MaxStack() uint {
	return self.maxStack
}
func (self *Method) MaxLocals() uint {
	return self.maxLocals
}
func (self *Method) Code() []byte {
	return self.code
}
func (self *Method) ArgSlotCount() uint {
	return self.argSlotCount
}

/**
调用ExceptionTable.findExceptionHandler（）方法搜索异常处理表，如果
能够找到对应的异常处理项，则返回它的handlerPc字段，否则返回–1。
**/
func (self *Method) FindExceptionHandler(exClass *Class, pc int) int {
	handler := self.exceptionTable.findExceptionHandler(exClass, pc)
	if handler != nil {
		return handler.handlerPc
	}
	return -1
}

// 获取方法对应的行号
/**
和源文件名一样，并不是每个方法都有行号表。如果方法没有
行号表，自然也就查不到pc对应的行号，这种情况下返回–1。本地
方法没有字节码，这种情况下返回–2。
**/
func (self *Method) GetLineNumber(pc int) int {
	if self.IsNative() {
		return -2
	}
	if self.lineNumberTable == nil {
		return -1
	}
	return self.lineNumberTable.GetLineNumber(pc)
}
