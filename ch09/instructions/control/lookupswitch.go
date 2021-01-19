package control

import "jvmgo/ch09/instructions/base"
import "jvmgo/ch09/rtda"

// 控制指令共有11条。jsr和ret指令在Java 6之前用于实现finally子句，
// 从Java 6开始，Oracle的Java编译器已经不再使用这两条指令了，
// 这里不讨论这两条指令。return系列指令有6条，用于从方法调用中返回。
// 剩下的3条指令为：goto、tableswitch和lookupswitch。


// 控制指令共有11条。jsr和ret指令在Java 6之前用于实现finally子句，
// 从Java 6开始，Oracle的Java编译器已经不再使用这两条指令了，
// 这里不讨论这两条指令。return系列指令有6条，用于从方法调用中返回。
// 剩下的3条指令为：goto、tableswitch和lookupswitch。

/**
Java语言中的switch-case语句有两种实现方式：如果case值可以编码成一个索引表，
则实现成tableswitch指令；否则实现成lookupswitch指令。
Java虚拟机规范的3.10小节里有两个例子，可以借用一下。



下面这个Java方法中的switch-case可以编译成tableswitch指令，代码如下：


int chooseNear(int i) {
    switch (i) {
        case 0: return 0;
        case 1: return 1;
        case 2: return 2;
        default: return -1;
    }
}




下面这个Java方法中的switch-case则需要编译成lookupswitch指令：

int chooseFar(int i) {
    switch (i) {
        case -100: return -1;
        case 0: return 0;
        case 100: return 1;
        default: return -1;
    }
}
**/



/*
lookupswitch
<0-3 byte pad>
defaultbyte1
defaultbyte2
defaultbyte3
defaultbyte4
npairs1
npairs2
npairs3
npairs4
match-offset pairs...
*/
// Access jump table by key match and jump


/**
lookupswitch中的matchOffsets有点像Map，它的key是case值，value是跳转偏移
量。Execute（）方法先从操作数栈中弹出一个int变量，然后用它查找
matchOffsets，看是否能找到匹配的key。如果能，则按照value给出的
偏移量跳转，否则按照defaultOffset跳转。
**/
type LOOKUP_SWITCH struct {
	defaultOffset int32
	npairs        int32
	matchOffsets  []int32
}

func (self *LOOKUP_SWITCH) FetchOperands(reader *base.BytecodeReader) {
	reader.SkipPadding()
	self.defaultOffset = reader.ReadInt32()
	self.npairs = reader.ReadInt32()
	self.matchOffsets = reader.ReadInt32s(self.npairs * 2)
}

func (self *LOOKUP_SWITCH) Execute(frame *rtda.Frame) {
	key := frame.OperandStack().PopInt()
	for i := int32(0); i < self.npairs*2; i += 2 {
		if self.matchOffsets[i] == key {
			offset := self.matchOffsets[i+1]
			base.Branch(frame, int(offset))
			return
		}
	}
	base.Branch(frame, int(self.defaultOffset))
}
