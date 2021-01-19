package classfile

/**
按照用途，23种预定义属性可以分为3组。第一组属性是实现Java虚拟机所必需的，共有5种；
第二组属性是Java类库所必需的，共有12种；
第三组属性主要是提供给工具使用，共有6种。第三组属性是可选的，也就是说可以不出现在class文件中。如果class文件中存在第三组属性，
Java虚拟机实现或者Java类库也是可以利用它们的，比如使用LineNumberTable属性在异常堆栈中显式行号。



从class文件演进的角度来讲，JDK1.0时只有6种预定义属性，JDK1.1增加了3种。J2SE 5.0增加了9种属性，
主要用于支持泛型和注解。Java SE 6增加了StackMapTable属性，用于优化字节码验证。Java SE 7增加了BootstrapMethods属性，
用于支持新增的invokedynamic指令。Java SE 8又增加了三种属性。



| 属性名                               | Java SE | 分组 | 位置                                  |
| ------------------------------------ | ------- | ---- | ------------------------------------- |
| ConstantValue                        | 1.0.2   | 1    | field_info                            |
| Code                                 | 1.0.2   | 1    | method_info                           |
| Exceptions                           | 1.0.2   | 1    | method_info                           |
| SourceFile                           | 1.0.2   | 3    | ClassFile                             |
| LineNumberTable                      | 1.0.2   | 3    | Code                                  |
| LocalVariableTable                   | 1.0.2   | 3    | Code                                  |
| InnerClasses                         | 1.1     | 2    | ClassFile                             |
| Synthetic                            | 1.1     | 2    | ClassFile,field_info,method_info      |
| Deprecated                           | 1.1     | 3    | ClassFile,field_info,method_info      |
| EnclosingMethod                      | 5.0     | 2    | ClassFile                             |
| Signature                            | 5.0     | 2    | ClassFile,field_info,method_info      |
| SourceDebugExtension                 | 5.0     | 3    | ClassFile                             |
| LocalVariableTypeTable               | 5.0     | 3    | Code                                  |
| RuntimeVisibleAnnotations            | 5.0     | 2    | ClassFile,field_info,method_info      |
| RuntimeInvisibleAnnotations          | 5.0     | 2    | ClassFile,field_info,method_info      |
| RuntimeVisibleParameterAnnotations   | 5.0     | 2    | method_info                           |
| RuntimeInvisibleParameterAnnotations | 5.0     | 2    | method_info                           |
| AnnotationDefault                    | 5.0     | 2    | method_info                           |
| StackMapTable                        | 6       | 1    | Code                                  |
| BoostrapMethods                      | 7       | 1    | ClassFile                             |
| RuntimeVisibleTypeAnnotations        | 8       | 2    | ClassFile,field_info,method_info,Code |
| RuntimeInvisibleTypeAnnotations      | 8       | 2    | ClassFile,field_info,method_info,Code |
| MethodParameters                     | 8       | 2    | method_info                           |
**/
/*
attribute_info {
    u2 attribute_name_index;
    u4 attribute_length;
    u1 info[attribute_length];
}
*/
type UnparsedAttribute struct {
	name   string
	length uint32
	info   []byte
}

func (self *UnparsedAttribute) readInfo(reader *ClassReader) {
	self.info = reader.readBytes(self.length)
}

func (self *UnparsedAttribute) Info() []byte {
	return self.info
}
