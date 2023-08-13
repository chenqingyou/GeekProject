package main

import "fmt"

func main() {
	var sliceInt = []int{1, 2, 3, 4, 5, 6}
	//将3删除
	fmt.Println(deleteSlice(6, sliceInt))
}

// 传入要删除的下标，返回一个切片
func deleteSlice(index int, sliceInt []int) []int {
	if index < len(sliceInt) {
		firstSlice := sliceInt[0:index]
		endSlice := sliceInt[index+1:]
		var resultSlice []int
		resultSlice = append(resultSlice, firstSlice...)
		resultSlice = append(resultSlice, endSlice...)
		return resultSlice
	}
	return sliceInt
}

// 优化后的版本
func deleteSliceNew(index int, sliceInt []any) []any {
	if index < len(sliceInt) || index < 0 {
		//直接修改原来的切片
		return append(sliceInt[:index], sliceInt[index+1:]...)
	}
	return sliceInt
}

// 传入要增加的下标，返回一个切片
func appendSlice(index int, value any, sliceInt []any) []any {
	if index < len(sliceInt)-1 || index < 0 {
		var resultSlice []any
		resultSlice = append(resultSlice, sliceInt[:index])
		resultSlice = append(resultSlice, value)
		resultSlice = append(resultSlice, sliceInt[index+1:]...)
		return resultSlice

	}
	return append(sliceInt, value)
}
