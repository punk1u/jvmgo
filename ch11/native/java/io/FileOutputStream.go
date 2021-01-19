package io

import "os"
import "unsafe"
import "jvmgo/ch11/native"
import "jvmgo/ch11/rtda"

const fos = "java/io/FileOutputStream"

func init() {
	native.Register(fos, "writeBytes", "([BIIZ)V", writeBytes)
}

// 表示FileOutputStream.writeBytes（）
// private native void writeBytes(byte b[], int off, int len, boolean append) throws IOException;
// ([BIIZ)V
func writeBytes(frame *rtda.Frame) {
	vars := frame.LocalVars()
	//this := vars.GetRef(0)
	b := vars.GetRef(1)
	off := vars.GetInt(2)
	len := vars.GetInt(3)
	//append := vars.GetBoolean(4)

	jBytes := b.Data().([]int8)
	/**
	虽然同是字节类型，但是在Java语言中byte是有符号类型，在
	Go语言中byte则是无符号类型。所以这里需要先把Java的字节数组
	转换成Go的[]byte变量，然后再调用os.Stdout.Write（）方法把它写到
	控制台。
	**/
	goBytes := castInt8sToUint8s(jBytes)
	goBytes = goBytes[off : off+len]
	os.Stdout.Write(goBytes)
}

// 将有符号类型转换成无符号类型
func castInt8sToUint8s(jBytes []int8) (goBytes []byte) {
	ptr := unsafe.Pointer(&jBytes)
	goBytes = *((*[]byte)(ptr))
	return
}
