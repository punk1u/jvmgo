<center><b>指令集和解释器</b></center>

# 字节码和指令集

每一个类或者接口都会被Java编译器编译成一个class文件，类或接口的方法信息就放在class文件的method_info结构中 。如果方法不是抽象的，也不是本地方法，方法的Java代码就会被编译器编译成字节码（即使方法是空的，编译器也会生成一条return语句），存放在method_info结构的Code属性中。



字节码中存放编码后的Java虚拟机指令。每条指令都以一个单字节的操作码（opcode）开头，这就是字节码名称的由来。由于只使用一字节表示操作码，显而易见，Java虚拟机最多只能支持256条指令。到第八版为止，Java虚拟机规范已经定义了205条指令，操作码分别是0（0x00）到202（0xCA）254（0xFE）和255（0xFF）。这205条指令构成了Java虚拟机的指令集（instruction set）。和汇编语言类似，为了便于记忆，Java虚拟机规范给每个操作码都指定了一个助记符（mnemonic）。比如操作码是0x00这条指令，因为它什么也不做，所以它的助记符是nop（no operation）。



Java虚拟机使用的是变长指令，操作码后面可以跟零字节或多字节的操作数（operand）。如果把指令想象成函数的话，操作数就是它的参数。为了让编码后的字节码更加紧凑，很多操作码本身就隐含了操作数，比如把常数0推入操作数栈的指令是iconst_0。



***操作数栈和局部变量表只存放数据的值，并不记录数据类型。结果就是：指令必须知道自己在操作什么类型的数据。这一点也直接反映在了操作码的助记符上。例如，iadd指令就是对int值进行加法操作；dstore指令把操作数栈顶的double值弹出，存储到局部变量表中；areturn从方法中返回引用值。也就是说，如果某类指令可以操作不同类型的变量，则助记符的第一个字母表示变量类型。***助记符首字母和变量类型的对应关系如下所示：



| 助记符首字母 | 数据类型     | 例子                   |
| ------------ | ------------ | ---------------------- |
| a            | reference    | aload、astore、areturn |
| b            | byte/boolean | bipush、baload         |
| c            | char         | caload、castore        |
| d            | double       | dload、dstore、dadd    |
| f            | float        | fload、fstore、fadd    |
| i            | int          | iload、istore、iadd    |
| l            | long         | load、lstore、ladd     |
| s            | short        | sipush、sastore        |



Java虚拟机规范把已经定义的205条指令按用途分成了11类，分别是：常量（constants）指令、加载（loads）指令、存储（stores）指令、操作数栈（stack）指令、数学（math）指令、转换（conversions）指令、比较（comparisons）指令、控制（control）指令、引用（references）指令、扩展（extended）指令和保留（reserved）指令。



保留指令一共有3条。其中一条是留给调试器的，用于实现断点，操作码是202（0xCA），助记符是breakpoint。另外两条留给Java虚拟机实现内部使用，操作码分别是254（0xFE）和266（0xFF），助记符是impdep1和impdep2。这三条指令不允许出现在class文件中。



Java虚拟机规范的2.11节介绍了Java虚拟机解释器的大致逻辑，如下所示：

```java
do {
    atomically calculate pc and fetch opcode at pc;
    if (operands) fetch operands;
    execute the action for the opcode;
} while (there is more to do);
```



# 指令

## 常量指令

常量指令把常量推入操作数栈顶。常量可以来自三个地方：隐含在操作码里、操作数和运行时常量池。常量指令共有21条。



### nop指令

nop指令是最简单的一条指令，因为它什么也不做。



### const系列指令

这一系列指令把隐含在操作码中的常量值推入操作数栈顶。



### bipush和sipush指令

bipush指令从操作数中获取一个byte型整数，扩展成int型，然后推入栈顶。sipush指令从操作数中获取一个short型整数，扩展成int型，然后推入栈顶。



## 加载指令

加载指令从局部变量表获取变量，然后推入操作数栈顶。加载指令共33条，按照所操作变量的类型可以分为6类：aload系列指令操作引用类型变量、dload系列操作double类型变量、fload系列操作float变量、iload系列操作int变量、lload系列操作long变量、xaload操作数组。



## 存储指令

和加载指令刚好相反，存储指令把变量从操作数栈顶弹出，然后存入局部变量表。和加载指令一样，存储指令也可以分为6类。



## 栈指令

栈指令直接对操作数栈进行操作，共9条：pop和pop2指令将栈顶变量弹出，dup系列指令复制栈顶变量，swap指令交换栈顶的两个变量。



和其他类型的指令不同，栈指令并不关心变量类型。





## 数学指令

数学指令大致对应Java语言中的加、减、乘、除等数学运算符。数学指令包括算术指令、位移指令和布尔运算指令等，共37条。



### 算术指令

算术指令又可以进一步分为加法（add）指令、减法（sub）指令、乘法（mul）指令、除法（div）指令、求余（rem）指令和取反（neg）指令6种。加、减、乘、除和取反指令都比较简单。



### 位移指令

位移指令可以分为左移和右移两种，右移指令又可以分为算术右移（有符号右移）和逻辑右移（无符号右移）两种。算术右移和逻辑位移的区别仅在于符号位的扩展，如下面的Java代码所示。

```java
int x = -1;
println(Integer.toBinaryString(x)); // 11111111111111111111111111111111
println(Integer.toBinaryString(x >> 8)); // 11111111111111111111111111111111
println(Integer.toBinaryString(x >>> 8)); // 00000000111111111111111111111111
```



### 布尔运算指令

布尔运算指令只能操作int和long变量，分为按位与（and）、按位或（or）、按位异或（xor）3种。



### iinc指令

iinc指令给局部变量表中的int变量增加常量值，局部变量表索引和常量值都由指令的操作数提供。



## 类型转换指令

类型转换指令大致对应Java语言中的基本类型强制转换操作。类型转换指令有共15条，引用类型转换对应的是checkcast指令。



按照被转换变量的类型，类型转换指令可以分为3种：i2x系列指令把int变量强制转换成其他类型；l2x系列指令把long变量强制转换成其他类型；f2x系列指令把float变量强制转换成其他类型；d2x系列指令把double变量强制转换成其他类型。



## 比较指令

比较指令可以分为两类：一类将比较结果推入操作数栈顶，一类根据比较结果跳转。比较指令是编译器实现if-else、for、while等语句的基石，共有19条。



### lcmp指令

lcmp指令用于比较long变量。



### fcmp和dcmp指令

fcmpg和fcmpl指令用于比较float变量。这两条指令和lcmp指令很像，但是除了比较的变量类型不同以外，还有一个重要的区别。由于浮点数计算有可能产生NaN（Not a Number）值，所以比较两个浮点数时，除了大于、等于、小于之外，还有第4种结果：无法比较。fcmpg和fcmpl指令的区别就在于对第4种结果的定义。



当两个float变量中至少有一个是NaN时，用fcmpg指令比较的结果是1，而用fcmpl指令比较的结果是-1。



### dcmpg和dcmpl指令

dcmpg和dcmpl指令用来比较double变量，这两条指令和fcmpg、fcmpl指令除了比较的变量类型不同之外并无区别。



### if指令

if<cond>指令把操作数栈顶的int变量弹出，然后跟0进行比较，满足条件则跳转。

假设从栈顶弹出的变量是x，则指令执行跳转操作的条件如下：

1. ifeq：x==0
2. ifne：x！=0
3. iflt：x<0
4. ifle：x<=0
5. ifgt：x>0
6. ifge：x>=0



### if_icmp<cond>指令

if_icmp<cond>指令把栈顶的两个int变量弹出，然后进行比较，满足条件则跳转。跳转条件和if<cond>指令类似。





### if_acmp<cond>指令

if_acmpeq和if_acmpne指令把栈顶的两个引用弹出，根据引用是否相同进行跳转。



## 控制指令

控制指令共有11条。jsr和ret指令在Java 6之前用于实现finally子句，从Java 6开始，Oracle的Java编译器已经不再使用这两条指令了，这里不讨论这两条指令。return系列指令有6条，用于从方法调用中返回。剩下的3条指令为：goto、tableswitch和lookupswitch。



### goto指令

goto指令进行无条件跳转。



### tableswitch指令

Java语言中的switch-case语句有两种实现方式：如果case值可以编码成一个索引表，则实现成tableswitch指令；否则实现成lookupswitch指令。Java虚拟机规范的3.10小节里有两个例子，可以借用一下。



下面这个Java方法中的switch-case可以编译成tableswitch指令，代码如下：

```java
int chooseNear(int i) {
    switch (i) {
        case 0: return 0;
        case 1: return 1;
        case 2: return 2;
        default: return -1;
    }
}
```



下面这个Java方法中的switch-case则需要编译成lookupswitch指令：

```java
int chooseFar(int i) {
    switch (i) {
        case -100: return -1;
        case 0: return 0;
        case 100: return 1;
        default: return -1;
    }
}
```





tableswitch指令操作码的后面有0~3字节的padding，以保证defaultOffset在字节码中的地址是4的倍数。 

```
tableswitch
<0-3 byte pad>
defaultbyte
lowbyte
highbyte
jump offsets...
```



### lookupswitch指令

lookupswitch结构：

```
<0-3 byte pad>
defaultbyte
npairs-num
npairs1
npairs2
npairs3
npairs4
match-offset pairs...
```



lookupswitch中的matchOffsets有点像Map，它的key是case值，value是跳转偏移量。Execute（）方法先从操作数栈中弹出一个int变量，然后用它查找
matchOffsets，看是否能找到匹配的key。如果能，则按照value给出的偏移量跳转，否则按照defaultOffset跳转。



## 扩展指令

扩展指令共有6条。



### wide指令

加载类指令、存储类指令、ret指令和iinc指令需要按索引访问局部变量表，索引以uint8的形式存在字节码中。对于大部分方法来说，局部变量表大小都不会超过256，所以用一字节来表示索引就够了。但是如果有方法的局部变量表超过这限制呢？Java虚拟机规范定义了wide指令来扩展前述指令。



### ifnull和ifnonnull指令

根据引用是否是null进行跳转，ifnull和ifnonnull指令把栈顶的引用弹出。



### goto_w指令

goto_w指令和goto指令的唯一区别就是索引从2字节变成了4字节。



