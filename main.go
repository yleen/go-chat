/*
程序主入口
*/
package main

func main() {
	// 两个程序都是main包 所以不用import 详见https://learnku.com/go/t/32464
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}

//执行方式
//go build -o main.go server.go
//./server
// nc 127.0.0.1 8888
