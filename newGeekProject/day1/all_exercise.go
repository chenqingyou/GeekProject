package main

import "fmt"

func sliceAdd(intSlice []int) {
	fmt.Printf("%p\n", intSlice)
	intSlice = append(intSlice, 4)
	for i := range intSlice {
		intSlice[i]++
	}
}

func main() {
	s1 := []int{1, 2, 3}
	s2 := s1
	s2 = append(s1, 0)
	sliceAdd(s1)
	sliceAdd(s2)
	fmt.Printf("s1:[%v][%p],s2:[%v][%p]", s1, s1, s2, s2)
	/*
			这是因为 Go 语言的 append 函数的特性。
			在 Go 语言中，当切片需要增长（通过 append 操作）并且其底层的数组无法存储更多的元素（也就是说，与底层数组的容量相比，长度已经达到了最大值）时
		    ，Go 语言会重新为这个切片分配一片新的内存然后将所有元素复制到这个新的内存中。新的切片将包含所有的元素，并且其底层数组的容量会增加。
			在你的代码中，当你将 s1 赋值给 s2 时，s1 和 s2 分享同一个底层数组。执行 s2 = append(s1, 0) 语句后，由于 s1 的容量已经满了，
		    系统会为 s2 的增长创建新的内存空间，所以 s1 和 s2 从此开始有了各自的内存空间，任何对 s2 的修改不会影响到 s1。
			当你调用 sliceAdd(s1) 和 sliceAdd(s2) 时，当 append 函数试图在 s1 的切片上添加元素时，由于满足扩容条件，这会创建一个新的底层数组（无法预测它的地址），
		    然而在函数返回时，这个新的切片并没有返回给 main 函数中的 s1。因此，在 main 函数中，s1 还是原来的切片。
			而对于 s2 在 main 函数中已经是自己独立的一份数据了，sliceAdd(s2) 对 s2 的修改会影响到 main 函数中的 s2。
			这就是为什么 s1 在 sliceAdd 之后没有变化，而 s2 有变化的原因。
	*/

}
