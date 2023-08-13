package main

import "fmt"

//切片的扩容原理：重新分配一段连续的内存，而后把原本的数据拷贝过去

func main() {
	s1 := []int{1, 2, 3, 4}
	fmt.Printf("s1:%v,len:%v,cap:%v\n", s1, len(s1), cap(s1))
	s1 = append(s1, 5) //扩容
	fmt.Printf("s1:%v,len:%v,cap:%v\n", s1, len(s1), cap(s1))
	/*
		1、当容量小于256时候，两倍扩容；否则安装1.25倍扩容
	*/
}
