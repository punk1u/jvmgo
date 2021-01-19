package classfile


// 表示Java虚拟机规范中的CONSTANT_NameAndType_info类型
//CONSTANT_NameAndType_info给出字段或方法的名称和描述符。CONSTANT_Class_info和CONSTANT_NameAndType_info加在一起可以唯一确定一个字段或者方法。

/**
字段名或方法名由name_index给出，字段或方法的描述符由descriptor_index给出。name_index和descriptor_index都是常量池索引，
指向CONSTANT_Utf8_info常量。字段和方法名就是代码中出现的（或者编译器生成的）字段或方法的名字。
Java虚拟机规范定义了一种简单的语法来描述字段和方法，可以根据下面的规则生成描述符。

1. 类型描述符
   1. 基本类型byte、short、char、int、long、float和double的描述符是单个字母，分别对应B、S、C、I、J、F和D。注意，long的描述符是J而不是L。
   2. 引用类型的描述符是L+类的完全限定名+分号
   3. 数组类型的描述符是[+数组元素类型描述符
2. 字段描述符就是字段类型的描述符
3. 方法描述符是（分号分隔的参数类型描述符）+返回值类型描述符，其中void返回值由单个字母v表示
| 字段描述符          | 字段类型           | 方法描述符            | 方法                                |
| ------------------- | ------------------ | --------------------- | ----------------------------------- |
| S                   | short              | ()v                   | void run()                          |
| L.java.lang.Object; | java.lang.Object   | ()Ljava.lang.String;  | String toString()                   |
| [I                  | int[]              | (Ljava.lang.String;)V | void main(String[] args)            |
| [[D                 | double[] []        | (FF)F                 | int max(float x,foat y)             |
| [Ljava.lang.String; | java.lang.Object[] | ([JJ)I                | int binarySearch(long[] a,long key) |

Java语言支持方法重载（override），不同的方法可以有相同的名字，只要参数列表不同即可。
这就是为什么CONSTANT_NameAndType_info结构要同时包含名称和描述符的原因。Java中不能定义多个同名字段，哪怕它们的类型各不相同，
但这只是Java语法的限制而已，从class文件的层面来看，是完全可以支持这点的。
**/
/*
CONSTANT_NameAndType_info {
    u1 tag;
    u2 name_index;
    u2 descriptor_index;
}
*/
type ConstantNameAndTypeInfo struct {
	nameIndex       uint16
	descriptorIndex uint16
}

func (self *ConstantNameAndTypeInfo) readInfo(reader *ClassReader) {
	self.nameIndex = reader.readUint16()
	self.descriptorIndex = reader.readUint16()
}
