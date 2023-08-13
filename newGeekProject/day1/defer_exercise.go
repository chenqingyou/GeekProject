package main

import "fmt"

// 函数中defer数量少于等于8个，第九个可执行可不执行
// 函数中defer关键字不能再循环中执行,编译的时候不知道有几个defer
// 函数中defer语句和return语句的乘机小于或者等于15个
// defer在栈上会比堆更加快，一个栈上只能有一个
func main() {
	DeferV1()
	DeferV2()
	DeferV3()
}

func DeferV1() {
	for i := 0; i < 10; i++ {
		defer func() { //传入的是一个地址，在循环中i的地址是不会变化的
			fmt.Printf("%v ", i)
		}()
	}
}

func DeferV2() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			fmt.Printf("%v ", j)
		}()
	}
}

func DeferV3() {
	for i := 0; i < 10; i++ {
		defer func(value int) {
			fmt.Printf("%v ", value)
		}(i)
	}
}
