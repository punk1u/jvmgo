package extended

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/instructions/loads"
import "jvmgo/ch11/instructions/math"
import "jvmgo/ch11/instructions/stores"
import "jvmgo/ch11/rtda"

// 加载类指令、存储类指令、ret指令和iinc指令需要按索引访问局部变量表，
// 索引以uint8的形式存在字节码中。对于大部分方法来说，局部变量表大小都不会超过256，
// 所以用一字节来表示索引就够了。但是如果有方法的局部变量表超过这限制呢？
// Java虚拟机规范定义了wide指令来扩展前述指令。

// wide指令改变其他指令的行为，modifiedInstruction字段存放被
// 改变的指令。wide指令需要自己解码出modifiedInstruction

// Extend local variable index by additional bytes
type WIDE struct {
	modifiedInstruction base.Instruction
}

// FetchOperands（）方法先从字节码中读取一字节的操作码，然后创建子指令实例，最后读取子指令的操作数。
// 加载指令和存储指令都只有一个操作数，需要扩展成2字节
func (self *WIDE) FetchOperands(reader *base.BytecodeReader) {
	opcode := reader.ReadUint8()
	switch opcode {
	case 0x15:
		inst := &loads.ILOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x16:
		inst := &loads.LLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x17:
		inst := &loads.FLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x18:
		inst := &loads.DLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x19:
		inst := &loads.ALOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x36:
		inst := &stores.ISTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x37:
		inst := &stores.LSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x38:
		inst := &stores.FSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x39:
		inst := &stores.DSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x3a:
		inst := &stores.ASTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x84:
		inst := &math.IINC{}
		inst.Index = uint(reader.ReadUint16())
		inst.Const = int32(reader.ReadInt16())
		self.modifiedInstruction = inst
	// 因为没有实现ret指令，所以暂时调用panic（）函数终止程序执行。
	case 0xa9: // ret
		panic("Unsupported opcode: 0xa9!")
	}
}

// wide指令只是增加了索引宽度，并不改变子指令操作，所以其Execute（）方法只要调用子指令的Execute（）方法即可
func (self *WIDE) Execute(frame *rtda.Frame) {
	self.modifiedInstruction.Execute(frame)
}
