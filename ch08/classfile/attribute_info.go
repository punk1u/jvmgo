package classfile

var (
	_attrDeprecated = &DeprecatedAttribute{}
	_attrSynthetic  = &SyntheticAttribute{}
)

/**
和常量池类似，各种属性表达的信息也各不相同，因此无法用统一的结构来定义。
不同之处在于，常量是由Java虚拟机规范严格定义的，共有14种。但属性是可扩展的，
不同的虚拟机实现可以定义自己的属性类型。由于这个原因，Java虚拟机规范没有使用tag，
而是使用属性名来区别不同的属性。属性数据放在属性名之后的u1表中，这样Java虚拟机实现就可以跳过自己无法识别的属性。
**/
/*
attribute_info {
    u2 attribute_name_index;
    u4 attribute_length;
    u1 info[attribute_length];
}
*/
type AttributeInfo interface {
	readInfo(reader *ClassReader)
}

// 读取属性表
func readAttributes(reader *ClassReader, cp ConstantPool) []AttributeInfo {
	attributesCount := reader.readUint16()
	attributes := make([]AttributeInfo, attributesCount)
	for i := range attributes {
		attributes[i] = readAttribute(reader, cp)
	}
	return attributes
}

// 读取单个属性，先读取属性名索引，根据它从常量池中找到属性名，然后读取属性长度，接着调用newAttributeInfo()函数以创建具体的属性实例
func readAttribute(reader *ClassReader, cp ConstantPool) AttributeInfo {
	attrNameIndex := reader.readUint16()
	attrName := cp.getUtf8(attrNameIndex)
	attrLen := reader.readUint32()
	attrInfo := newAttributeInfo(attrName, attrLen, cp)
	attrInfo.readInfo(reader)
	return attrInfo
}

func newAttributeInfo(attrName string, attrLen uint32, cp ConstantPool) AttributeInfo {
	switch attrName {
	// case "AnnotationDefault":
	case "BootstrapMethods":
		return &BootstrapMethodsAttribute{}
	case "Code":
		return &CodeAttribute{cp: cp}
	case "ConstantValue":
		return &ConstantValueAttribute{}
	case "Deprecated":
		return _attrDeprecated
	case "EnclosingMethod":
		return &EnclosingMethodAttribute{cp: cp}
	case "Exceptions":
		return &ExceptionsAttribute{}
	case "InnerClasses":
		return &InnerClassesAttribute{}
	case "LineNumberTable":
		return &LineNumberTableAttribute{}
	case "LocalVariableTable":
		return &LocalVariableTableAttribute{}
	case "LocalVariableTypeTable":
		return &LocalVariableTypeTableAttribute{}
	// case "MethodParameters":
	// case "RuntimeInvisibleAnnotations":
	// case "RuntimeInvisibleParameterAnnotations":
	// case "RuntimeInvisibleTypeAnnotations":
	// case "RuntimeVisibleAnnotations":
	// case "RuntimeVisibleParameterAnnotations":
	// case "RuntimeVisibleTypeAnnotations":
	case "Signature":
		return &SignatureAttribute{cp: cp}
	case "SourceFile":
		return &SourceFileAttribute{cp: cp}
	// case "SourceDebugExtension":
	// case "StackMapTable":
	case "Synthetic":
		return _attrSynthetic
	default:
		return &UnparsedAttribute{attrName, attrLen, nil}
	}
}
