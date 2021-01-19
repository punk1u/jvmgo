package classfile

import "fmt"
import "unicode/utf16"

/*
CONSTANT_Utf8_info {
    u1 tag;
    u2 length;
    u1 bytes[length];
}
*/
// 表示虚拟机规范中的CONSTANT_Utf8_info
type ConstantUtf8Info struct {
	str string
}

// 先读取出[]byte，然后调用decodeMUTF8()函数把它解码成Go字符串
func (self *ConstantUtf8Info) readInfo(reader *ClassReader) {
	length := uint32(reader.readUint16())
	bytes := reader.readBytes(length)
	self.str = decodeMUTF8(bytes)
}

// 获取对应的字符串数据
func (self *ConstantUtf8Info) Str() string {
	return self.str
}

/**
注意：字符串在class文件中是以MUTF-8（Modified UTF-8）方式编码的，而不是标准的UTF-8编码格式。MUTF-8编码方式和UTF-8大致相同，但并不兼容。差别有两点：

1. null字符（代码点U+0000）会被编码为2字节：0xC0、0x80
2. 补充字符（Supplementary Characters，代码点大于U+FFFF的Unicode字符）是按UTF-16拆分为代理对（Surrogate Pair）分别编码的。
**/
// 因为Go语言字符串使用UTF-8编码，所以如果字符串中不包含null字符或补充字符，下面这个简化版的readMUTF8()也可以正常工作
/*
func decodeMUTF8(bytes []byte) string {
	return string(bytes) // not correct!
}
*/

// mutf8 -> utf16 -> utf32 -> string
// see java.io.DataInputStream.readUTF(DataInput)
func decodeMUTF8(bytearr []byte) string {
	utflen := len(bytearr)
	chararr := make([]uint16, utflen)

	var c, char2, char3 uint16
	count := 0
	chararr_count := 0

	for count < utflen {
		c = uint16(bytearr[count])
		if c > 127 {
			break
		}
		count++
		chararr[chararr_count] = c
		chararr_count++
	}

	for count < utflen {
		c = uint16(bytearr[count])
		switch c >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			/* 0xxxxxxx*/
			count++
			chararr[chararr_count] = c
			chararr_count++
		case 12, 13:
			/* 110x xxxx   10xx xxxx*/
			count += 2
			if count > utflen {
				panic("malformed input: partial character at end")
			}
			char2 = uint16(bytearr[count-1])
			if char2&0xC0 != 0x80 {
				panic(fmt.Errorf("malformed input around byte %v", count))
			}
			chararr[chararr_count] = c&0x1F<<6 | char2&0x3F
			chararr_count++
		case 14:
			/* 1110 xxxx  10xx xxxx  10xx xxxx*/
			count += 3
			if count > utflen {
				panic("malformed input: partial character at end")
			}
			char2 = uint16(bytearr[count-2])
			char3 = uint16(bytearr[count-1])
			if char2&0xC0 != 0x80 || char3&0xC0 != 0x80 {
				panic(fmt.Errorf("malformed input around byte %v", (count - 1)))
			}
			chararr[chararr_count] = c&0x0F<<12 | char2&0x3F<<6 | char3&0x3F<<0
			chararr_count++
		default:
			/* 10xx xxxx,  1111 xxxx */
			panic(fmt.Errorf("malformed input around byte %v", count))
		}
	}
	// The number of chars produced may be less than utflen
	chararr = chararr[0:chararr_count]
	runes := utf16.Decode(chararr)
	return string(runes)
}
