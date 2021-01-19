package classfile


/**
表示Code变长属性，只存在于method_info结构中。Code属性中存放字节码等方法相关信息。
max_stack给出操作数栈的最大深度，max_locals给出局部变量表大小。接着是字节码，存在u1表中。最后是异常处理表和属性表。
在第4章讨论运行时数据区，并且实现操作数栈和局部变量表时，max_stack和max_locals就会派上用场。
在第5章讨论指令集和解释器时，会用到字节码。在第10章讨论异常处理时，会使用异常处理表。
**/
/*
Code_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 max_stack;
    u2 max_locals;
    u4 code_length;
    u1 code[code_length];
    u2 exception_table_length;
    {   u2 start_pc;
        u2 end_pc;
        u2 handler_pc;
        u2 catch_type;
    } exception_table[exception_table_length];
    u2 attributes_count;
    attribute_info attributes[attributes_count];
}
*/
type CodeAttribute struct {
	cp             ConstantPool
	maxStack       uint16
	maxLocals      uint16
	code           []byte
	exceptionTable []*ExceptionTableEntry
	attributes     []AttributeInfo
}

// 将ClassReader中的数据读取到Code属性中对应的值中
func (self *CodeAttribute) readInfo(reader *ClassReader) {
	self.maxStack = reader.readUint16()
	self.maxLocals = reader.readUint16()
	codeLength := reader.readUint32()
	self.code = reader.readBytes(codeLength)
	self.exceptionTable = readExceptionTable(reader)
	self.attributes = readAttributes(reader, self.cp)
}

func (self *CodeAttribute) MaxStack() uint {
	return uint(self.maxStack)
}
func (self *CodeAttribute) MaxLocals() uint {
	return uint(self.maxLocals)
}
func (self *CodeAttribute) Code() []byte {
	return self.code
}
func (self *CodeAttribute) ExceptionTable() []*ExceptionTableEntry {
	return self.exceptionTable
}

type ExceptionTableEntry struct {
	startPc   uint16
	endPc     uint16
	handlerPc uint16
	catchType uint16
}

// 读取异常处理表的数据
func readExceptionTable(reader *ClassReader) []*ExceptionTableEntry {
	exceptionTableLength := reader.readUint16()
	exceptionTable := make([]*ExceptionTableEntry, exceptionTableLength)
	for i := range exceptionTable {
		exceptionTable[i] = &ExceptionTableEntry{
			startPc:   reader.readUint16(),
			endPc:     reader.readUint16(),
			handlerPc: reader.readUint16(),
			catchType: reader.readUint16(),
		}
	}
	return exceptionTable
}

func (self *ExceptionTableEntry) StartPc() uint16 {
	return self.startPc
}
func (self *ExceptionTableEntry) EndPc() uint16 {
	return self.endPc
}
func (self *ExceptionTableEntry) HandlerPc() uint16 {
	return self.handlerPc
}
func (self *ExceptionTableEntry) CatchType() uint16 {
	return self.catchType
}
