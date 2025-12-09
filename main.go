package main

import "go-admin-full/cmd/godash"

// main 函数是程序的入口点
// 该函数不接受任何参数，也不返回任何值
// 主要功能是调用 godash.Execute() 函数来执行核心逻辑
func main() {
	godash.Execute()
}
