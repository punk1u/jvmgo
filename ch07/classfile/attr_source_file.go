package classfile

// 表示SourceFile属性，用于指出源文件名
/*
SourceFile_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 sourcefile_index;
}
*/
type SourceFileAttribute struct {
	cp              ConstantPool
	sourceFileIndex uint16
}

// 获取源文件名对应的在常量池中的索引
func (self *SourceFileAttribute) readInfo(reader *ClassReader) {
	self.sourceFileIndex = reader.readUint16()
}

// 获取源文件名
func (self *SourceFileAttribute) FileName() string {
	return self.cp.getUtf8(self.sourceFileIndex)
}
