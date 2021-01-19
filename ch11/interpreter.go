package main

import "fmt"
import "jvmgo/ch11/instructions"
import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"

func interpret(thread *rtda.Thread, logInst bool) {
	defer catchErr(thread)
	loop(thread, logInst)
}

// 把局部变量表和操作数栈的内容打印出来，以此来观察方法的执行结果。
func catchErr(thread *rtda.Thread) {
	if r := recover(); r != nil {
		logFrames(thread)
		panic(r)
	}
}

/**
在每次循环开始，先拿到当前帧，然后根据pc从当前方法中解
码出一条指令。指令执行完毕之后，判断Java虚拟机栈中是否还有
帧。如果没有则退出循环；否则继续。
**/
func loop(thread *rtda.Thread, logInst bool) {
	reader := &base.BytecodeReader{}
	for {
		frame := thread.CurrentFrame()
		pc := frame.NextPC()
		thread.SetPC(pc)

		// decode
		reader.Reset(frame.Method().Code(), pc)
		opcode := reader.ReadUint8()
		inst := instructions.NewInstruction(opcode)
		inst.FetchOperands(reader)
		frame.SetNextPC(reader.PC())

		if logInst {
			logInstruction(frame, inst)
		}

		// execute
		inst.Execute(frame)
		if thread.IsStackEmpty() {
			break
		}
	}
}

// logInstruction（）函数在方法执行过程中打印指令信息
func logInstruction(frame *rtda.Frame, inst base.Instruction) {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	pc := frame.Thread().PC()
	fmt.Printf("%v.%v() #%2d %T %v\n", className, methodName, pc, inst, inst)
}

// logFrames（）函数打印Java虚拟机栈信息
func logFrames(thread *rtda.Thread) {
	for !thread.IsStackEmpty() {
		frame := thread.PopFrame()
		method := frame.Method()
		className := method.Class().Name()
		lineNum := method.GetLineNumber(frame.NextPC())
		fmt.Printf(">> line:%4d pc:%4d %v.%v%v \n",
			lineNum, frame.NextPC(), className, method.Name(), method.Descriptor())
	}
}
