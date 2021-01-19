<center><b>类加载器</b></center>



# 类加载器


## System类是如何被初始化的

System类有3个公开的静态常量：out、err和in。其中out和err用于向标准输出流和标准错误流中写入信息，in用于从标准输入流中读取信息。

System类的源代码：

```java
// java.lang.System
public final class System {
    public final static InputStream in = null;
    public final static PrintStream out = null;
    public final static PrintStream err = null;
    /* register the natives via the static initializer.
    *
    * VM will invoke the initializeSystemClass method to complete
    * the initialization for this class separated from clinit.
    * Note that to use properties set by the VM, see the constraints
    * described in the initializeSystemClass method.
    */
    private static native void registerNatives();
    static {
    	registerNatives();
    }
    // 其他代码
}
```

从注释可知，System类的初始化过程分为两个阶段。第一个阶段由类初始化方法完成，在这个方法中registerNatives（）方法会注册其他本地方法。第二个阶段由VM完成，在这个阶段VM会调用System.initializeSystemClass（）方法。

那么initializeSystemClass（）方法究竟干了些什么呢？这个方法很长，而且有很详细的注释。去掉无关的代码和注释之后，它的代码如下：

```java
/**
* Initialize the system class. Called after thread initialization.
*/
private static void initializeSystemClass() {
     // 其他代码
    FileInputStream fdIn = new FileInputStream(FileDescriptor.in);
    FileOutputStream fdOut = new FileOutputStream(FileDescriptor.out);
    FileOutputStream fdErr = new FileOutputStream(FileDescriptor.err);
    setIn0(new BufferedInputStream(fdIn));
    setOut0(newPrintStream(fdOut, props.getProperty("sun.stdout.encoding")));
    setErr0(newPrintStream(fdErr, props.getProperty("sun.stderr.encoding")));
    // 其他代码
}
```

可见in、out和err正是在这里设置的。再来看sun.misc.VM类的源代码（VM类属于Oracle私有代码，并没有开源，下面是反编译后的Java代码）:

```java
// sun.misc.VM
public class VM {
    // 其他代码
    static {
    // 其他代码
    	initialize();
    }
    private static native void initialize();
}
```

VM类在初始化时调用了initialize（）方法。虽然initialize（）是本地方法，但是可以推测正是这个方法调用了System.initializeSystemClass（）方法。



## System.out.println（）是如何工作的

回到System.initializeSystemClass（）方法，进一步省略之后，其代码如下：

```java
// java.lang.System
public final static PrintStream out = null;
private static void initializeSystemClass() {
    // 其他代码
    FileOutputStream fdOut = new FileOutputStream(FileDescriptor.out);
    setOut0(newPrintStream(fdOut, props.getProperty("sun.stdout.encoding")));
    // 其他代码
}
```



setOut0（）是个本地方法，代码如下：

```java
private static native void setOut0(PrintStream out);
```



newPrintStream（）方法的代码如下：

```java
private static PrintStream newPrintStream(FileOutputStream fos, String enc) {
    if (enc != null) {
    try {
    	return new PrintStream(new BufferedOutputStream(fos, 128), true, enc);
    } catch (UnsupportedEncodingException uee) {}
    }
    return new PrintStream(new BufferedOutputStream(fos, 128), true);
}
```

由代码可知，System.out常量是PrintStream类型，它内部包装了一个BufferedOutputStream实例。BufferedOutputStream内部又包装了一个FileOutputStream实例。Java的io类库使用了装饰器模式，调用System.out.println（String）方法之后，经过层层包装，最后到达FileOutputStream类的writeBytes（）方法。这个方法无法用Java代码实现，所以是个本地方法，其代码如下：

```java
// java.io.FileOutputStream
public class FileOutputStream extends OutputStream {
    // 其他代码
    private native void writeBytes(byte b[], int off, int len, boolean append)
    throws IOException;
}
```

