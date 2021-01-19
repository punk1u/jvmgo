<center><b>本地方法调用</b></center>



# 本地方法调用

要想运行Java程序，除了Java虚拟机之外，还需要Java类库的配合。Java虚拟机和Java类库一起构成了Java运行时环境。Java类库主要用Java语言编写，一些无法用Java语言实现的方法则使用本地语言编写，这些方法叫作本地方法。



## 注册和查找本地方法

本地方法注册表，用来注册和查找本地方法。



Java虚拟机规范并没有规定如何实现和调用本地方法。



## 反射

### 类和对象之间的关系

在Java中，类也表现为普通的对象，它的类是java.lang.Class。

Java有强大的反射能力。可以在运行期间获取类的各种信息、存取静态和实例变量、调用方法，等等。要想运用这种能力，获取类对象是第一步。在Java语言中，有两种方式可以获得类对象引用：使用类字面值和调用对象的getClass（）方法。

下面的Java代码演示了这两种方式。

```java
System.out.println(String.class);
System.out.println("abc".getClass());
```



### 基本类型的类

void和基本类型也有对应的类对象，但只能通过字面值来访问，如下面的Java代码所示。

```java
System.out.println(void.class);
System.out.println(boolean.class);
System.out.println(byte.class);
System.out.println(char.class);
System.out.println(short.class);
System.out.println(int.class);
System.out.println(long.class);
System.out.println(float.class);
System.out.println(double.class);
```

和数组类一样，基本类型的类也是由Java虚拟机在运行期间生成的。



1. void和基本类型的类名就是void、int、float等。
2. 基本类型的类没有超类，也没有实现任何接口。
3. 非基本类型的类对象是通过ldc指令加载到操作数栈中的



而基本类型的类对象，虽然在Java代码中看起来是通过字面量获取的，但是编译之后的指令并不是ldc，而是getstatic。每个基本类型都有一个包装类，包装类中有一个静态常量，叫作TYPE，其中存放的就是基本类型的类。例如
java.lang.Integer类，代码如下：

```java
public final class Integer extends Number implements Comparable<Integer> {
    // 其他代码
    @SuppressWarnings("unchecked")
    public static final Class<Integer> TYPE=(Class<Integer>)Class.getPrimitiveClass("int");
    // 其他代码
}
```

也就是说，基本类型的类是通过getstatic指令访问相应包装类的TYPE字段加载到操作数栈中的。



### 通过反射获取类名

为了支持通过反射获取类名，需要实现以下4个本地方法：

1. java.lang.Object.getClass()
2. java.lang.Class.getPrimitiveClass()
3. java.lang.Class.getName0()
4. java.lang.Class.desiredAssertionStatus0()

Object.getClass（）就不用多说了，它返回对象的类对象引用。Class.getPrimitiveClass()之前提到过，基本类型的包装类在初始化时会调用这个方法给TPYE字段赋值。Character类是基本类型char的包装类，它在初始化时会调用Class.desiredAssertionStatus0（）方法，所以这个方法也需要实现。最后，之所以要实现getName0（）方法，是因为Class.getName()方法是依赖这个本地方法工作的，该方法的代码如下：

```java
// java.lang.Class
public String getName() {
    String name = this.name;
    if (name == null)
    	this.name = name = getName0();
    return name;
}
```



## 字符串拼接和String.intern()方法

### Java类库

在Java语言中，通过加号来拼接字符串。作为优化，javac编译器会把字符串拼接操作转换成StringBuilder的使用。



例如下面这段Java代码：

```java
String hello = "hello,";
String world = "world!";
String str = hello + world;
System.out.println(str);
```

很可能会被javac优化为下面这样：

```java
String str = new StringBuilder().append("hello,").append("world!").toString();
System.out.println(str);
```

为了运行上面的代码，本节将实现以下3个本地方法：

1. System.arrayCopy（）
2. Float.floatToRawIntBits（）
3. Double.doubleToRawLongBits（）



String.getChars（）方法调用了System.arraycopy（）方法拷贝数组，代码如下：

```java
// java.lang.String
public void getChars(int srcBegin, int srcEnd, char dst[], int dstBegin) {
     // 其他代码
     System.arraycopy(value, srcBegin, dst, dstBegin, srcEnd - srcBegin);
}
```

StringBuilder最终会调用到Math类，Math类在初始化时需要调用Float.floatToRawIntBits（）和
Double.doubleToRawLongBits（）方法，代码如下：

```java
package java.lang;
public final class Math {
    // Use raw bit-wise conversions on guaranteed non-NaN arguments.
    private static long negativeZeroFloatBits = Float.floatToRawIntBits(-0.0f);
    private static long negativeZeroDoubleBits = Double.doubleToRawLongBits(-0.0d);
}
```

Float.floatToRawIntBits（）和Double.doubleToRawLongBits（）返回浮点数的编码。



## Object.hashCode()、equals()和toString()

Object类有3个非常重要的方法：hashCode（）返回对象的哈希码；equals（）用来比较两个对象是否“相同”；toString（）返回对象的字符串表示。hashCode（）是个本地方法，equals（）和toString（）则是用Java写的，它们的代码如下：

```java
package java.lang;
public class Object {
    // 其他代码省略
    public native int hashCode();
    public boolean equals(Object obj) {
    	return (this == obj);
    }
    public String toString() {
    	return getClass().getName() + "@" + Integer.toHexString(hashCode());
    }
}
```



## Object.clone()

Object类提供了clone（）方法，用来支持对象克隆。这也是一个本地方法，代码如下：

```java
// java.lang.Object
protected native Object clone() throws CloneNotSupportedException;
```

数组也实现了Cloneable接口，所以需要处理数组克隆的情况。

## 自动装箱和拆箱

为了更好地融入Java的对象系统，每种基本类型都有一个包装类与之对应。从Java 5开始，Java语法增加了自动装箱和拆箱（autoboxing/unboxing）能力，可以在必要时把基本类型转换成包装类型或者反之。这个增强完全是由编译器完成的，Java虚拟机没有做任何调整。



以int类型为例，它的包装类是java.lang.Integer。它提供了2个方法来帮助编译器在int变量和Integer对象之间转换：静态方法valueOf()把int变量包装成Integer对象；实例方法intValue（）返回被包装的int变量。这两个方法的代码如下：

```java
package java.lang;
public final class Integer extends Number implements Comparable<Integer> {
    // 其他代码省略
    private final int value;
    public static Integer valueOf(int i) {
        if (i >= IntegerCache.low && i <= IntegerCache.high)
        return IntegerCache.cache[i + (-IntegerCache.low)];
        return new Integer(i);
    }
    public int intValue() {
    	return value;
    }
}
```



由上面的代码可知，Integer.valueOf（）方法并不是每次都创建Integer（）对象，而是维护了一个缓存池IntegerCache。对于比较小（默认是–128~127）的int变量，在IntegerCache初始化时就预先加载到了池中，需要用时直接从池里取即可。IntegerCache是Integer类的内部类，为了便于参考，下面给出它的完整代码。

```java
private static class IntegerCache {
    static final int low = -128;
    static final int high;
    static final Integer cache[];
    static {
        int h = 127; // high value may be configured by property
        String integerCacheHighPropValue =
        	sun.misc.VM.getSavedProperty("java.lang.Integer.IntegerCache.high");
        if (integerCacheHighPropValue != null) {
            try {
                int i = parseInt(integerCacheHighPropValue);
                i = Math.max(i, 127);
                // Maximum array size is Integer.MAX_VALUE
                h = Math.min(i, Integer.MAX_VALUE - (-low) -1);
            } catch( NumberFormatException nfe) {
            	// If the property cannot be parsed into an int, ignore it.
            }
        }
        high = h;
        cache = new Integer[(high - low) + 1];
        int j = low;
        for(int k = 0; k < cache.length; k++)
        	cache[k] = new Integer(j++);
        // range [-128, 127] must be interned (JLS7 5.1.7)
        assert IntegerCache.high >= 127;
    }
    private IntegerCache() {}
}
```

需要说明的是IntegerCache在初始化时需要确定缓存池中Integer对象的上限值，为此它调用了sun.misc.VM类的getSavedProperty()方法。



