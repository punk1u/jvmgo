package classfile

import "encoding/binary"

// 作用：帮助读取class中的数据
type ClassReader struct {
	data []byte
}

// 读取u1类型数据,读取data的第一个字节的数据，并将data的值改变为其原来第一个字节之后的数据
func (self *ClassReader) readUint8() uint8 {
	val := self.data[0]
	self.data = self.data[1:]
	return val
}

// 读取u2类型数据
// Go标准库encoding/binary包中定义了一个变量BigEndian，正好可以从[]byte中解码多字节数据
func (self *ClassReader) readUint16() uint16 {
	val := binary.BigEndian.Uint16(self.data)
	self.data = self.data[2:]
	return val
}

// 读取u4类型数据
func (self *ClassReader) readUint32() uint32 {
	val := binary.BigEndian.Uint32(self.data)
	self.data = self.data[4:]
	return val
}

// 读取u8int64（JAVA虚拟机规范并没有定义u8）类型数据
func (self *ClassReader) readUint64() uint64 {
	val := binary.BigEndian.Uint64(self.data)
	self.data = self.data[8:]
	return val
}

// 读取uint16表，表的大小由开头的uint16数据指出
func (self *ClassReader) readUint16s() []uint16 {
	n := self.readUint16()
	s := make([]uint16, n)
	for i := range s {
		s[i] = self.readUint16()
	}
	return s
}

// 用于读取指定数量的字节
func (self *ClassReader) readBytes(n uint32) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}
