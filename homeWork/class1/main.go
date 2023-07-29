package main

import "fmt"

func main() {
	var sliceInt = []int{1, 2, 3, 4, 5, 6}
	//将3删除
	fmt.Println(deleteSlice(2, sliceInt))
	fmt.Println(deleteSliceAny(2, sliceInt)) //会改变原切片
	fmt.Printf("sliceInt cap[%v]\n", cap(sliceInt))
	var sliceInts1 = []int{1, 2, 3, 4, 5, 6}
	sliceInt1 := deleteAndShrink(2, sliceInts1) //删除对应的下班，并且对切片进行缩容
	fmt.Printf("结果[%v],缩容后的sliceInt1 cap[%v]\n", sliceInt1, cap(sliceInt))
}

// 传入要删除的下标，返回一个切片
func deleteSlice(index int, sliceInt []int) []int {
	if index < len(sliceInt) {
		firstSlice := sliceInt[0:index] //首
		endSlice := sliceInt[index+1:]  //尾
		var resultSlice []int
		resultSlice = append(resultSlice, firstSlice...)
		resultSlice = append(resultSlice, endSlice...)
		return resultSlice
	}
	return sliceInt
}

// 使用泛型优化后的版本
func deleteSliceAny[T any](index int, sliceInt []T) []T {
	if index < len(sliceInt) || index < 0 {
		//直接修改原来的切片
		return append(sliceInt[:index], sliceInt[index+1:]...)
	}
	return sliceInt
}

// 传入要增加的下标，返回一个切片
func appendSlice[T any](index int, value T, sliceAny []T) []T {
	if index < len(sliceAny)-1 || index < 0 {
		var resultSlice []T
		resultSlice = append(resultSlice, sliceAny[:index]...)
		resultSlice = append(resultSlice, value)
		resultSlice = append(resultSlice, sliceAny[index+1:]...)
		return resultSlice
	}
	return append(sliceAny, value)
}

// 缩容
func deleteAndShrink[T any](index int, slice []T) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}
	newSlice := make([]T, len(slice)-1) //创建一个新的切片，容量是后面删除后的切片的长度
	copy(newSlice, slice[:index])       //Copy是将后面的拷贝到前面的切片中
	copy(newSlice[index:], slice[index+1:])
	return newSlice
}
