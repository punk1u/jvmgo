package main

import "fmt"
import "jvmgo/ch06/instructions"
import "jvmgo/ch06/instructions/base"
import "jvmgo/ch06/rtda"
import "jvmgo/ch06/rtda/heap"

func interpret(method *heap.Method) {
	// interpret（）方法的其余代码先创建一个Thread实例，然后创建一个帧并把它推入Java虚拟机栈顶，最后执行方法
	thread := rtda.NewThread()
	frame := thread.NewFrame(method)
	thread.PushFrame(frame)

	defer catchErr(frame)
	loop(thread, method.Code())
}

// 把局部变量表和操作数栈的内容打印出来，以此来观察方法的执行结果。
func catchErr(frame *rtda.Frame) {
	if r := recover(); r != nil {
		//fmt.Printf("LocalVars:%v\n", frame.LocalVars())
		//fmt.Printf("OperandStack:%v\n", frame.OperandStack())
		//panic(r)
	}
}

func loop(thread *rtda.Thread, bytecode []byte) {
	frame := thread.PopFrame()
	reader := &base.BytecodeReader{}

	for {
		pc := frame.NextPC()
		thread.SetPC(pc)

		// decode
		reader.Reset(bytecode, pc)
		opcode := reader.ReadUint8()
		inst := instructions.NewInstruction(opcode)
		inst.FetchOperands(reader)
		frame.SetNextPC(reader.PC())

		// execute
		fmt.Printf("pc:%2d inst:%T %v\n", pc, inst, inst)
		inst.Execute(frame)
	}
}
