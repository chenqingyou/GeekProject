package main

import "fmt"

func prepend[T Number](s []T, x T) []T {
	s = append(s, x) // 首先在切片尾部添加元素
	copy([]T{2}, s)  // copy函数的使用方式：首个参数应该是一个目标切片，第二个参数是源切片
	fmt.Println("s", s, copy([]T{2}, s))
	s[0] = x // 最后在切片头部位置插入元素
	return s
}

type Number interface {
	int | int64 | float64 //代表对泛型的约束
}

// Sum 使用泛型实现一个求和的参数  如果使用any，里面有一些不能被计算的类型
func Sum[T Number](vals ...T) T {
	var res T
	for _, val := range vals {
		res = res + val
	}
	return res
}

func main() {
	x := []int{1, 2, 3}
	fmt.Println(prepend(x, 0)) // Output: [0 1 2 3]
}
