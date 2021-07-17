/*
服务端基本构建
*/
package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	//创建一个server对象
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	//因为上面赋值的是& 即返回的就是指针
	return server
}

func (this *Server) Handler(conn net.Conn) {
	//...当前链接的业务
	fmt.Println("链接建立成功")
}

//启动服务器接口
func (this *Server) Start() {
	//socker listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net listen err:", err)
		return
	}

	//close listen socket
	defer listener.Close()
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}

}
