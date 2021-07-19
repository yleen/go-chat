package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个UserApi
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听当前user channel 消息的goroutine
	go user.ListenMessage()
	return user
}

//用户上线业务
func (this *User) Online() {
	//用户上线，将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//广播当前用户消息
	this.server.BroadCast(this, "已上线")
}

//用户下线业务
func (this *User) Offline() {
	//用户下线，将用户加入到onlineMap中
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播当前用户消息
	this.server.BroadCast(this, "已下线")
}

//用户接收消息
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg + "\n"))
}

//用户处理消息业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Name + "]" + "--" + user.Addr + ":" + "在线...\n"
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("该用户名已存在")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMessage("用户名" + newName + "设置成功")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//私聊消息，格式：to|name｜msg
		//1 获取用户名
		toUser := strings.Split(msg, "|")[1]
		if toUser == "" {
			this.SendMessage("请输入私聊的用户名\n")
			return
		}
		//2 根据用户名 得到对方User对象
		sendUser, ok := this.server.OnlineMap[toUser]
		if !ok {
			this.SendMessage("该用户不存在\n")
			return
		}
		//3 获取消息内容 通过对方的User对象将消息内容发送过去
		toMsg := strings.Split(msg, "|")[2]
		if toMsg == "" {
			this.SendMessage("消息内容为空")
			return
		}
		sendUser.SendMessage(this.Name + "发起私聊:" + toMsg)
	} else {
		this.server.BroadCast(this, msg)
	}
}

//监听当前User channel的方法，一旦有消息直接发给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
