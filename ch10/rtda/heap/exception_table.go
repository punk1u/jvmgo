package heap

import "jvmgo/ch10/classfile"


type ExceptionTable []*ExceptionHandler

// 表示异常处理表
/**
异常处理表的每一项都包含3个信息：处理哪部分代码抛出的异常、哪类异常，
以及异常处理代码在哪里。具体来说，start_pc和end_pc可以锁定一部分字节码，
这部分字节码对应某个可能抛出异常的try{}代码块。catch_type是个索引，
通过它可以从运行时常量池中查到一个类符号引用，解析后的类是个异常类，
假定这个类是X。如果位于start_pc和end_pc之间的指令抛出异常x，
且x是X（或者X的子类）的实例，handler_pc就指出负责异常处理的catch{}块在哪里。
**/
type ExceptionHandler struct {
	startPc   int
	endPc     int
	handlerPc int
	catchType *ClassRef
}

/**
把class文件中的异常处理表转换成ExceptionTable类型。
有一点需要特别说明：异常处理项的catchType有可能是0。
我们知道0是无效的常量池索引，但是在这里0并非表示catch-none，
而是表示catch-all
**/
func newExceptionTable(entries []*classfile.ExceptionTableEntry, cp *ConstantPool) ExceptionTable {
	table := make([]*ExceptionHandler, len(entries))
	for i, entry := range entries {
		table[i] = &ExceptionHandler{
			startPc:   int(entry.StartPc()),
			endPc:     int(entry.EndPc()),
			handlerPc: int(entry.HandlerPc()),
			catchType: getCatchType(uint(entry.CatchType()), cp),
		}
	}

	return table
}

// 从运行时常量池中查找类符号引用
func getCatchType(index uint, cp *ConstantPool) *ClassRef {
	if index == 0 {
		return nil // catch all
	}
	return cp.GetConstant(index).(*ClassRef)
}

// 查找异常处理表
/**
异常处理表查找逻辑

void catchOne() {
    try {
    	tryItOut();
    } catch (TestExc e) {
    	handleExc(e);
    }
}

当tryItOut()方法通过athrow指令抛出TestExc异常时，
Java虚拟机首先会查找tryItOut()方法的异常处理表，
看它能否处理该异常。如果能，则跳转到相应的字节码开始异常处理。
假设tryItOut()方法无法处理异常，Java虚拟机会进一步查看它的调用者，
也就是catchOne()方法的异常处理表。catchOne()方法刚好可以处理TestExc异常，
使catch{}块得以执行。假设catchOne()方法也无法处理TestExc异常，
Java虚拟机会继续查找catchOne()的调用者的异常处理表。
这个过程会一直继续下去，直到找到某个异常处理项，或者到达Java虚拟机栈的底部。

这里注意两点。
第一，startPc给出的是try{}语句块的第一条指令，
endPc给出的则是try{}语句块的下一条指令。
第二，如果catchType是nil（在class文件中是0），
表示可以处理所有异常，这是用来实现finally子句的。
**/
func (self ExceptionTable) findExceptionHandler(exClass *Class, pc int) *ExceptionHandler {
	for _, handler := range self {
		// jvms: The start_pc is inclusive and end_pc is exclusive
		if pc >= handler.startPc && pc < handler.endPc {
			if handler.catchType == nil {
				return handler
			}
			catchClass := handler.catchType.ResolvedClass()
			if catchClass == exClass || catchClass.IsSuperClassOf(exClass) {
				return handler
			}
		}
	}
	return nil
}
