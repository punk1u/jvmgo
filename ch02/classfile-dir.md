<center><b>类路径</b></center>



# 类路径

Java虚拟机规范并没有规定虚拟机应该从哪里寻找类，因此不同的虚拟机实现可以采用不同的方法。Oracle的Java虚拟机实现根据类路径（class path）来搜索类。按照搜索的先后顺序，类路径可以分为以下3个部分：

1. 启动类路径（bootstrap classpath）
2. 扩展类路径（extension classpath）
3. 用户类路径（user classpath）

启动类路径默认对应jre\lib目录，Java标准库（大部分在rt.jar里）位于该路径。扩展类路径默认对应jre\lib\ext目录，使用Java扩展机制的类位于这个路径。我们自己实现的类，以及第三方类库则位于用户类路径。可以通过-Xbootclasspath选项修改启动类路径。

Java虚拟机将使用JDK的启动类路径来寻找和加载Java标准库中的类，因此需要某种方式指定jre目录的位置。