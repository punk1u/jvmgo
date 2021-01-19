package heap

import "unicode/utf16"

/**
在class文件中，字符串是以MUTF8格式保存的，
在Java虚拟机运行期间，字符串以java.lang.String（后面简称String）对象的形式存在，
而在String对象内部，字符串又是以UTF16格式保存的。
字符串相关功能大部分都是由String（和StringBuilder等）类提供的。

String类有两个实例变量。其中一个是value，类型是字符数组，
用于存放UTF16编码后的字符序列。另一个是hash，缓存计字符串的哈希码。

字符串对象是不可变（immutable）的，一旦构造好之后，
就无法再改变其状态（这里指value字段）。

为了节约内存，Java虚拟机内部维护了一个字符串池。
String类提供了intern（）实例方法，可以把自己放入字符串池。
**/

// key:go string  value:java string
var internedStrings = map[string]*Object{}


/**
根据Go字符串返回相应的Java字符串实例。如果
Java字符串已经在池中，直接返回即可，否则先把Go字符串（UTF8
格式）转换成Java字符数组（UTF16格式），然后创建一个Java字符串
实例，把它的value变量设置成刚刚转换而来的字符数组，最后把
Java字符串放入池中。
**/
// todo
// go string -> java.lang.String
func JString(loader *ClassLoader, goStr string) *Object {
	if internedStr, ok := internedStrings[goStr]; ok {
		return internedStr
	}

	chars := stringToUtf16(goStr)
	jChars := &Object{loader.LoadClass("[C"), chars}

	jStr := loader.LoadClass("java/lang/String").NewObject()
	jStr.SetRefVar("value", "[C", jChars)

	internedStrings[goStr] = jStr
	return jStr
}

// java.lang.String -> go string
func GoString(jStr *Object) string {
	charArr := jStr.GetRefVar("value", "[C")
	return utf16ToString(charArr.Chars())
}

// Go语言字符串在内存中是UTF8编码的，先把它强制转成
// UTF32，然后调用utf16包的Encode（）函数编码成UTF16。
// utf8 -> utf16
func stringToUtf16(s string) []uint16 {
	runes := []rune(s)         // utf32
	return utf16.Encode(runes) // func Encode(s []rune) []uint16
}

// utf16 -> utf8
func utf16ToString(s []uint16) string {
	runes := utf16.Decode(s) // func Decode(s []uint16) []rune
	return string(runes)
}
