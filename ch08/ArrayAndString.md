<center><b>数组和字符串</b></center>

# 数组和字符串

在大部分编程语言中，数组和字符串都是最基本的数据类型。Java虚拟机直接支持数组，对字符串的支持则由java.lang.String和相关的类提供。



## 数组概述

数组在Java虚拟机中是个比较特殊的概念。有下面几个原因：

1. 首先，数组类和普通的类是不同的。普通的类从class文件中加载，但是数组类由Java虚拟机在运行时生成。数组的类名是左方括号（[）+数组元素的类型描述符；数组的类型描述符就是类名本身。例如，int[]的类名是[I，	int[] []的类名是[[I，Object[]的类名是[Ljava/lang/Object；，String[] []的类名是[[java/lang/String；，等等。
2. 其次，创建数组的方式和创建普通对象的方式不同。普通对象由new指令创建，然后由构造函数初始化。基本类型数组由newarray指令创建；引用类型数组由anewarray指令创建；另外还有一个专门的multianewarray指令用于创建多维数组。
3. 最后，很显然，数组和普通对象存放的数据也是不同的。普通对象中存放的是实例变量，通过putfield和getfield指令存取。数组对象中存放的则是数组元素，通过<t>aload和<t>astore系列指令按索引存取。其中<t>可以是a、b、c、d、f、i、l或者s，分别用于存取引用、byte、char、double、float、int、long或short类型的数组。另外，还有一个arraylength指令，用于获取数组长度。



### 数组对象

和普通对象一样，数组也是分配在堆中的，通过引用来使用。



数组类不需要初始化，数组类的超类是java.lang.Object，并且实现了java.lang.Cloneable和java.io.Serializable接口。



## 数组相关指令

### newarray指令

newarray指令用来创建基本类型数组，包括boolean[]、byte[]、char[]、short[]、int[]、long[]、float[]和double[]8种。



newarray指令需要两个操作数。第一个操作数是一个uint8整数，在字节码中紧跟在指令操作码后面，表示要创建哪种类型的数组。Java虚拟机规范把这个操作数叫作atype，并且规定了它的有效值。



newarray指令的第二个操作数是count，从操作数栈中弹出，表示数组长度。



### anewarray指令

anewarray指令用来创建引用类型数组。



anewarray指令也需要两个操作数。第一个操作数是uint16索引，来自字节码。通过这个索引可以从当前类的运行时常量池中找到一个类符号引用，解析这个符号引用就可以得到数组元素的类。第二个操作数是数组长度，从操作数栈中弹出。



如果是数组类名，描述符就是类名，直接返回即可。如果是基本类型名，返回对应的类型描述符，否则肯定是普通的类名，前面加上方括号，结尾加上分号即可得到类型描述符。



### arraylength指令

arraylength指令用于获取数组长度。arraylength指令只需要一个操作数，即从操作数栈顶弹出的数组引用。

arraylength指令只需要一个操作数，即从操作数栈顶弹出的数组引用。



###  < t > aload指令

 < t > aload系列指令按索引取数组元素值，然后推入操作数栈。



### < t > astore指令

< t >astore系列指令按索引给数组元素赋值。 



iastore指令的三个操作数分别是：要赋给数组元素的值、数组索引、数组引用，依次从操作数栈中弹出。如果数组引用是null，则抛出NullPointerException。如果数组索引小于0或者大于等于数组长度，则抛出ArrayIndexOutOfBoundsException异常。这两个检查和<t>aload系列指令一样。如果一切正常，则按索引给数组元素赋值。



### multianewarray指令

multianewarray指令创建多维数组。



multianewarray指令的第一个操作数是个uint16索引，通过这个索引可以从运行时常量池中找到一个类符号引用，解析这个引用就可以得到多维数组类。第二个操作数是个uint8整数，表示数组维度。这两个操作数在字节码中紧跟在指令操作码后面。multianewarray指令还需要从操作数栈中弹出n个整数，分别代表每一个维度的数组长度。



数组的类型检查和强制类型转换

1. 数组可以强制转换成Object类型（因为数组的超类是Object）。
2. 数组可以强制转换成Cloneable和Serializable类型（因为数组实现了这两个接口）。
3. 如果下面两个条件之一成立，类型为[]SC的数组可以强制转换成类型为[]TC的数组：
   1. TC和SC是同一个基本类型。
   2. TC和SC都是引用类型，且SC可以强制转换成TC。



## 字符串

在class文件中，字符串是以MUTF8格式保存的，在Java虚拟机运行期间，字符串以java.lang.String（后面简称String）对象的形式存在，而在String对象内部，字符串又是以UTF16格式保存的。字符串相关功能大部分都是由String（和StringBuilder等）类提供的。

String类有两个实例变量。其中一个是value，类型是字符数组，用于存放UTF16编码后的字符序列。另一个是hash，缓存计字符串的哈希码。

字符串对象是不可变（immutable）的，一旦构造好之后，就无法再改变其状态（这里指value字段）。

为了节约内存，Java虚拟机内部维护了一个字符串池。String类提供了intern（）实例方法，可以把自己放入字符串池。





