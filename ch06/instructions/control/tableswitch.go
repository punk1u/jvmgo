package control

import "jvmgo/ch06/instructions/base"
import "jvmgo/ch06/rtda"

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

/**
tableswitch指令操作码的后面有0~3字节的padding，以保证
defaultOffset在字节码中的地址是4的倍数。 

tableswitch
<0-3 byte pad>
defaultbyte1
defaultbyte2
defaultbyte3
defaultbyte4
lowbyte1
lowbyte2
lowbyte3
lowbyte4
highbyte1
highbyte2
highbyte3
highbyte4
jump offsets...
*/
// Access jump table by index and jump
type TABLE_SWITCH struct {
	defaultOffset int32
	low           int32
	high          int32
	jumpOffsets   []int32
}

func (self *TABLE_SWITCH) FetchOperands(reader *base.BytecodeReader) {
	reader.SkipPadding()
	self.defaultOffset = reader.ReadInt32()
	self.low = reader.ReadInt32()
	self.high = reader.ReadInt32()
	jumpOffsetsCount := self.high - self.low + 1
	self.jumpOffsets = reader.ReadInt32s(jumpOffsetsCount)
}

func (self *TABLE_SWITCH) Execute(frame *rtda.Frame) {
	index := frame.OperandStack().PopInt()

	var offset int
	if index >= self.low && index <= self.high {
		offset = int(self.jumpOffsets[index-self.low])
	} else {
		offset = int(self.defaultOffset)
	}

	base.Branch(frame, offset)
}
