<center><b>解析class文件</b></center>



# class文件



## class文件格式

作为类（或者接口） 信息的载体，每个class文件都完整地定义了一个类。为了使Java程序可以“编写一次，处处运行”，Java虚拟机规范对class文件格式进行了严格的规定。但是另一方面，对于从哪里加载class文件，给了足够多的自由。



Java虚拟机实现可以从文件系统读取和从JAR（或ZIP）压缩包中提取class文件。除此之外，也可以通过网络下载、从数据库加载，甚至是在运行中直接生成class文件。Java虚拟机规范（和本书）中所指的class文件，并非特指位于磁盘中的.class文件，而是泛指任何格式符合规范的class数据。



构成class文件的基本数据单位是字节，可以把整个class文件当成一个字节流来处理。稍大一些的数据由连续多个字节构成，这些数据在class文件中以大端（big-endian）方式存储。为了描述class文件格式，Java虚拟机规范定义了u1、u2和u4三种数据类型来表示1、2和4字节无符号整数。相同类型的多条数据一般按表（table）的形式存储在class文件中。表由表头和表项（item）构成，表头是u2或u4整数。假设表头是n，后面就紧跟着n个表项数据。



Java虚拟机规范使用一种类似C语言的结构体语法来描述class文件格式。整个class文件被描述为一个ClassFile结构，代码如下：

```java
ClassFile {
    u4 magic;
    u2 minor_version;
    u2 major_version;
    u2 constant_pool_count;
    cp_info constant_pool[constant_pool_count-1];
    u2 access_flags;
    u2 this_class;
    u2 super_class;
    u2 interfaces_count;
    u2 interfaces[interfaces_count];
    u2 fields_count;
    field_info fields[fields_count];
    u2 methods_count;
    method_info methods[methods_count];
    u2 attributes_count;
    attribute_info attributes[attributes_count];
}
```

获取类名可通过this_class的值作为索引从常量池中获取。SuperClassName同理。

###  魔数

很多文件格式都会规定满足该格式的文件必须以某几个固定字节开头，这几个字节主要起标识作用，叫作魔数（magic number）。例如PDF文件以4字节“%PDF”（0x25、0x50、0x44、0x46）开头，ZIP文件以2字节“PK”（0x50、0x4B）开头。class文件的魔数是“0xCAFEBABE”。

Java虚拟机规范规定，如果加载的class文件不符合要求的格式，Java虚拟机实现就抛出java.lang.ClassFormatError异常。



### 版本号

魔数之后是class文件的次版本号和主版本号，都是u2类型。假设某class文件的主版本号是M，次版本号是m，那么完整的版本号可以表示成“M.m”的形式。次版本号只在J2SE 1.2之前用过，从1.2开始基本上就没什么用了（都是0）。主版本号在J2SE 1.2之前是45，从1.2开始，每次有大的Java版本发布，都会加1。

特定的Java虚拟机实现只能支持版本号在某个范围内的class文件。Oracle的实现是完全向后兼容的，比如Java SE 8支持版本号为45.0~52.0的class文件。如果版本号不在支持的范围内，Java虚拟机实现就抛出java.lang.UnsupportedClassVersionError异常。



### 常量池

常量池占据了class文件很大一部分数据，里面存放着各式各样的常量信息，包括数字和字符串常量、类和接口名、字段和方法名，等等。

常量池实际上也是一个表，但是有三点需要特别注意。

1. ***表头给出的常量池大小比实际大1。假设表头给出的值是n，那么常量池的实际大小是n–1。***
2. ***有效的常量池索引是1~n–1。0是无效索引，表示不指向任何常量。***
3. ***CONSTANT_Long_info和CONSTANT_Double_info各占两个位置。也就是说，如果常量池中存在这两种常量，实际的常量数量比n–1还要少，而且1~n–1的某些数也会变成无效索引。***



由于常量池中存放的信息各不相同，所以每种常量的格式也不同。常量数据的第一字节是tag，用来区分常量类型。

Java虚拟机规范一共定义了14种常量:

| 常量类型           | JVM规范定义的常量类型对应的值 |
| ------------------ | ----------------------------- |
| Class              | 7                             |
| Fieldref           | 9                             |
| Methodref          | 10                            |
| InterfaceMethodref | 11                            |
| String             | 8                             |
| Integer            | 3                             |
| Float              | 4                             |
| Long               | 5                             |
| Double             | 6                             |
| NameAndType        | 12                            |
| Utf8               | 1                             |
| MethodHandle       | 15                            |
| MethodType         | 16                            |
| InvokeDynamic      | 18                            |



#### CONSTANT_Integer_info

CONSTANT_Integer_info使用4字节存储整数常量，其结构定义如下：

```
CONSTANT_Integer_info {
    u1 tag;
    u4 bytes;
}
```

CONSTANT_Integer_info正好可以容纳一个Java的int型常量，但实际上比int更小的boolean、byte、short和char类型的常量也放在CONSTANT_Integer_info中。



#### CONSTANT_Float_info

CONSTANT_Float_info使用4字节存储IEEE754单精度浮点数常量，结构如下：

```
CONSTANT_Float_info {
    u1 tag;
    u4 bytes;
}
```



#### CONSTANT_Long_info

CONSTANT_Long_info使用8字节存储整数常量，结构如下：

```
CONSTANT_Long_info {
    u1 tag;
    u4 high_bytes;
    u4 low_bytes;
}
```

这里使用了两个u4类型来表示CONSTANT_Long_info的8字节数据的高4位和低4位数据。



#### CONSTANT_Double_info

最后一个数字常量是CONSTANT_Double_info，使用8字节存储IEEE754双精度浮点数，结构如下：

```
CONSTANT_Double_info {
    u1 tag;
    u4 high_bytes;
    u4 low_bytes;
}
```



#### CONSTANT_Utf8_info

CONSTANT_Utf8_info常量里放的是MUTF-8编码的字符串，结构如下：

```
CONSTANT_Utf8_info {
    u1 tag;
    u2 length;
    u1 bytes[length];
}
```

字符串在class文件中是以MUTF-8（Modified UTF-8）方式编码的。而不是标准的UTF-8编码方式。

MUTF-8编码方式和UTF-8大致相同，但并不兼容。差别有两点：一是null字符（代码点U+0000）会被编码成2字节：0xC0、0x80；二是补充字符（Supplementary Characters，代码点大于U+FFFF的Unicode字符）是按UTF-16拆分为代理对（Surrogate Pair）分别编码的。



#### CONSTANT_String_info

CONSTANT_String_info常量表示java.lang.String字面量，结构如下：

```
CONSTANT_String_info {
    u1 tag;
    u2 string_index;
}
```

可以看到，CONSTANT_String_info本身并不存放字符串数据，只存了常量池索引，这个索引指向一个CONSTANT_Utf8_info常量。



#### CONSTANT_Class_info

CONSTANT_Class_info常量表示类或者接口的符号引用，结构如下：

```
CONSTANT_Class_info {
    u1 tag;
    u2 name_index;
}
```

和CONSTANT_String_info类似，name_index是常量池索引，指向CONSTANT_Utf8_info常量。



#### CONSTANT_NameAndType_info

CONSTANT_NameAndType_info给出字段或方法的名称和描述符。CONSTANT_Class_info和CONSTANT_NameAndType_info加在一起可以唯一确定一个字段或者方法。其结构如下：

```
CONSTANT_NameAndType_info {
    u1 tag;
    u2 name_index;
    u2 descriptor_index;
}
```

字段或方法名由name_index给出，字段或方法的描述符由descriptor_index给出。name_index和descriptor_index都是常量池索引，指向CONSTANT_Utf8_info常量。字段和方法名就是代码中出现的（或者编译器生成的）字段或方法的名字。

***Java虚拟机规范定义了一种简单的语法来描述字段和方法，可以根据下面的规则生成描述符。***

1. 类型描述符。
   1. 基本类型byte、short、char、int、long、float和double的描述符是单个字母，分别对应B、S、C、I、J、F和D。注意，long的描述符是J而不是L。
   2. 引用类型的描述符是L＋类的完全限定名＋分号。
   3. 数组类型的描述符是[＋数组元素类型描述符。
2. 字段描述符就是字段类型的描述符。
3. 方法描述符是（分号分隔的参数类型描述符）+返回值类型描述符，其中void返回值由单个字母V表示。

<center>字段和方法描述符示例</center>

| 字段类型         | 字段描述符          | 方法                                | 方法描述符             |
| ---------------- | ------------------- | ----------------------------------- | ---------------------- |
| short            | S                   | void run()                          | ()V                    |
| java.lang.Object | Ljava.lang.Object   | String toString()                   | ()Ljava.lang.String;   |
| [I               | int[]               | void main(String[] args)            | ([Ljava.lang.String;)V |
| double[] []      | [[D                 | int max(float x,float y)            | (FF)F                  |
| java.lang.Object | [Ljava.lang.Object; | int binarySearch(long[] a,long key) | ([JJ)I                 |

Java语言支持方法重载（override），不同的方法可以有相同的名字，只要参数列表不同即可。这就是为什么CONSTANT_NameAndType_info结构要同时包含名称和描述符的原因。那么字段呢？Java是不能定义多个同名字段的，哪怕它们的类型各不相同。这只是Java语法的限制而已，从class文件的层面来看，是完全可以支持这点的。



#### CONSTANT_Fieldref_info、CONSTANT_Methodref_info和CONSTANT_InterfaceMethodref_info
CONSTANT_Fieldref_info表示字段符号引用，CONSTANT_Methodref_info表示普通（非接口）方法符号引用，CONSTANT_InterfaceMethodref_info表示接口方法符号引用。这三种常量结构一模一样，为了节约篇幅，下面只给出CONSTANT_Fieldref_info的结构。

```
CONSTANT_Fieldref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
```

class_index和name_and_type_index都是常量池索引，分别指向CONSTANT_Class_info和CONSTANT_NameAndType_info常量。先定义一个统一的结构体ConstantMemberrefInfo来表示这3种常量。



#### CONSTANT_MethodType_info、CONSTANT_MethodHandle_info和CONSTANT_InvokeDynamic_info

CONSTANT_MethodType_info、CONSTANT_MethodHandle_info和CONSTANT_InvokeDynamic_info。它们是Java SE 7才添加到class文件中的，目的是支持新增的invokedynamic指令。





#### 常量池小结

可以把常量池中的常量分为两类：字面量（literal）和符号引用（symbolic reference）。字面量包括数字常量和字符串常量，符号引用包括类和接口名、字段和方法信息等。除了字面量，其他常量都是通过索引直接或间接指向CONSTANT_Utf8_info常量。





### 类访问标志

常量池之后是类访问标志，这是一个16位的“bitmask”，指出class文件定义的是类还是接口，访问级别是public还是private，等等。



### 类和超类索引

类访问标志之后是两个u2类型的常量池索引，分别给出类名和超类名。class文件存储的类名类似完全限定名，但是把点换成了斜线，Java语言规范把这种名字叫作二进制名（binary names）。因为每个类都有名字，所以thisClass必须是有效的常量池索引。除java.lang.Object之外，其他类都有超类，所以superClass只在Object.class中是0，在其他class文件中必须是有效的常量池索引。



### 接口索引表

类和超类索引后面是接口索引表，表中存放的也是常量池索引，给出该类实现的所有接口的名字。



### 字段和方法表

接口索引表之后是字段表和方法表，分别存储字段和方法信息。字段和方法的基本结构大致相同，差别仅在于属性表。下面是Java虚拟机规范给出的字段结构定义。

```
field_info {
    u2 access_flags;
    u2 name_index;
    u2 descriptor_index;
    u2 attributes_count;
    attribute_info attributes[attributes_count];
}
```

和类一样，字段和方法也有自己的访问标志。访问标志之后是一个常量池索引，给出字段名或方法名，然后又是一个常量池索引，给出字段或方法的描述符，最后是属性表。



### 属性表

方法的字节码便是存在属性表中。



和常量池类似，各种属性表达的信息也各不相同，因此无法用统一的结构来定义。不同之处在于，常量是由Java虚拟机规范严格定义的，共有14种。但属性是可以扩展的，不同的虚拟机实现可以定义自己的属性类型。由于这个原因，Java虚拟机规范没有使用tag，而是使用属性名来区别不同的属性。属性数据放在属性名之后的u1表中，这样Java虚拟机实现就可以跳过自己无法识别的属性。

```
attribute_info {
    u2 attribute_name_index;
    u4 attribute_length;
    u1 info[attribute_length];
}
```

注意，属性表中存放的属性名实际上并不是编码后的字符串，而是常量池索引，指向常量池中的CONSTANT_Utf8_info常量。



Java虚拟机规范预定义了23种属性，按照用途，23种预定义属性可以分为三组。

第一组属性是实现Java虚拟机所必需的，共有5种；

第二组属性是Java类库所必需的，共有12种；

第三组属性主要提供给工具使用，共有6种。第三组属性是可选的，也就是说可以不出现在class文件中。如果class文件中存在第三组属性，Java虚拟机实现或者Java类库也是可以利用它们的，比如使用LineNumberTable属性在异常堆栈中显示行号。



从class文件演进的角度来讲，JDK1.0时只有6种预定义属性，JDK1.1增加了3种。J2SE 5.0增加了9种属性，主要用于支持泛型和注解。Java SE 6增加了StackMapTable属性，用于优化字节码验证。Java SE 7增加了BootstrapMethods属性，用于支持新增的invokedynamic指令。Java SE 8又增加了三种属性。



下表给出了这23种属性出现的Java版本、分组以及它们在class文件中的位置。

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
| BootstrapMethods                     | 7       | 1    | ClassFile                             |
| RuntimeVisibleTypeAnnotations        | 8       | 2    | ClassFile,field_info,method_info,Code |
| RuntimeInvisibleTypeAnnotations      | 8       | 2    | ClassFile,field_info,method_info,Code |
| MethodParameters                     | 8       | 2    | method_info                           |



#### Deprecated和Synthetic属性

Deprecated和Synthetic是最简单的两种属性，仅起标记作用，不包含任何数据。这两种属性都是JDK1.1引入的，可以出现在ClassFile、field_info和method_info结构中，它们的结构定义如下：

```
Deprecated_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
}
Synthetic_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
}
```

由于不包含任何数据，所以attribute_length的值必须是0。



Deprecated属性用于指出类、接口、字段或方法已经不建议使用，编译器等工具可以根据Deprecated属性输出警告信息。

Synthetic属性用来标记源文件中不存在、由编译器生成的类成员，引入Synthetic属性主要是为了支持嵌套类和嵌套接口。





#### SourceFile属性

SourceFile是可选定长属性，用于指出源文件名。其结构定义如下：

```
SourceFile_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 sourcefile_index;
}
```

attribute_length的值必须是2。sourcefile_index是常量池索引，指向CONSTANT_Utf8_info常量。表示源文件名，如`ClassTest.java`，attribute_name_index则表示这个属性对应的属性名（SourceFile）在常量池中的索引。



#### ConstantValue属性

ConstantValue是定长属性，只会出现在field_info结构中，用于表示常量表达式的值（详见Java语言规范的15.28节）。其结构定义如下：

```
ConstantValue_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 constantvalue_index;
}
```

attribute_length的值必须是2。constantvalue_index是常量池索引，但具体指向哪种常量因字段类型而异。下表给出了字段类型和常量类型的对应关系。

| 字段类型                        | 常量类型              |
| ------------------------------- | --------------------- |
| long                            | CONSTANT_Long_info    |
| float                           | CONSTANT_Float_info   |
| double                          | CONSTANT_Double_info  |
| int、short、char、byte、boolean | CONSTANT_Integer_info |
| String                          | CONSTANT_String_info  |



#### Code属性

Code是变长属性，只存在于method_info结构中。Code属性中存放字节码等方法相关信息。相比前面介绍的几种属性，Code属性比较复杂，其结构定义如下：

```
Code_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 max_stack;
    u2 max_locals;
    u4 code_length;
    u1 code[code_length];
    u2 exception_table_length;
    { 
        u2 start_pc;
        u2 end_pc;
        u2 handler_pc;
        u2 catch_type;
    } 
    exception_table[exception_table_length];
    u2 attributes_count;
    attribute_info attributes[attributes_count];
}
```

max_stack给出操作数栈的最大深度，max_locals给出局部变量表大小。接着是字节码，存在u1表中。最后是异常处理表和属性表。



#### Exceptions属性

Exceptions是变长属性，记录方法抛出的异常表，其结构定义如下：

```
Exceptions_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 number_of_exceptions;
    u2 exception_index_table[number_of_exceptions];
}
```



#### LineNumberTable和LocalVariableTable属性

LineNumberTable属性表存放方法的行号信息，LocalVariableTable属性表中存放方法的局部变量信息。这两种属性和前面介绍的SourceFile属性都属于调试信息，都不是运行时必需的。在使用javac编译器编译Java程序时，默认会在class文件中生成这些信息。可以使用javac提供的-g：none选项来关闭这些信息的生成。



LineNumberTable和LocalVariableTable属性表在结构上很像，下面以LineNumberTable为例进行讨论，它的结构定义如下：

```
LineNumberTable_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 line_number_table_length;
    { 
        u2 start_pc;
        u2 line_number;
    } 
    line_number_table[line_number_table_length];
}
```



## Java语言基本数据类型



| JAVA语言类型 | 说明                 |
| ------------ | -------------------- |
| byte         | 8比特有符号整数      |
| short        | 16比特有符号整数     |
| char         | 16比特无符号整数     |
| int          | 32比特有符号整数     |
| long         | 64比特有符号整数     |
| float        | 32比特IEEE-754浮点数 |
| double       | 64比特IEEE-754浮点数 |



