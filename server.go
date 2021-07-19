/*
服务端基本构建
*/
package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	//创建一个server对象
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	//因为上面赋值的是& 即返回的就是指针
	return server
}

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		//将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//...当前链接的业务
	// fmt.Println("链接建立成功")
	user := NewUser(conn, this)

	user.Online()

	//监听用户是否活跃的channel
	isLive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			//io.EOF 表示文件流输入结束
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用户的消息 （去除‘\n’）
			msg := string(buf[:n-1])
			//用户针对msg进行消息处理
			user.DoMessage(msg)

			isLive <- true
		}
	}()

	//当前handler阻塞
	for {
		select {
		//超时强踢功能
		case <-isLive:
		//当前用户是活跃的，应该重置定时器
		//不做任何操作，为了激活select 更新定时器
		case <-time.After(time.Second * 1000):
			//已经超时
			//将当前User强制关闭
			user.SendMessage("连接超时")
			//销毁资源
			close(user.C)
			//退出当前Handler
			conn.Close()
			return
		}
	}

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

	//启动监听Message的goroutine
	go this.ListenMessager()
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
