package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>>链接服务器失败")
		return

	}
	//开启一个goroutine去处理server回执信息
	go client.DealResponse()

	fmt.Println(">>>>>>>链接服务器成功")

	//启动客户端的业务
	client.Run()
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址，默认是127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器")
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("链接server失败，err:", err)
		return nil
	}
	client.conn = conn
	//返回对象
	return client
}

//处理server回应的消息，直接显示到标准输出
func (client *Client) DealResponse() {
	//一旦client.conn有数据， 就直接copy到stdout标准输出上  永久阻塞监听
	io.Copy(os.Stdout, client.conn)

	//等价于
	/*
		for {
			buf := make()
			client.conn.Read(buf)
			fmt.Println(buf)
		}
	*/
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1 公聊模式")
	fmt.Println("2 私聊模式")
	fmt.Println("3 更新用户名")
	fmt.Println("0 退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag < 4 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入正确的数字")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式已选择")
			client.PublicChat()
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式已选择")
			client.PrivateChat()
			break
		case 3:
			//更新用户名
			fmt.Println("更新用户名已选择")
			client.UpdateName()
			break
		}
	}
}

//公聊方法
func (client *Client) PublicChat() {
	var msg string
	fmt.Println(">>>>>>>>>输入...按回车发送消息 exit 退出")
	fmt.Scanln(&msg)
	for msg != "exit" {
		if len(msg) != 0 {
			sendMsg := msg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("消息发送失败", err)
				break
			}
		}
		msg = ""
		fmt.Println(">>>>>>>>>输入...按回车发送消息 exit 退出")
		fmt.Scanln(&msg)
	}
}

//私聊方法
func (client *Client) PrivateChat() {

	var toUser string
	var msg string

	_, err := client.conn.Write([]byte("who\n"))
	if err != nil {
		fmt.Println("获取在线用户失败 请重试", err)
		return
	}
	fmt.Println(">>>>>>>>>输入您要私聊的用户名称 exit退出私聊模式")
	fmt.Scanln(&toUser)
	fmt.Println(">>>>>>>>>输入...按回车发送消息 exit退出")
	fmt.Scanln(&msg)
	for toUser != "exit" {
		for msg != "exit" {
			if len(msg) != 0 {
				sendMsg := "to|" + toUser + "|" + msg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("消息发送失败", err)
					break
				}
			}
			msg = ""
			fmt.Println(">>>>>>>>>输入...按回车发送消息 exit退出")
			fmt.Scanln(&msg)
		}
		fmt.Println(">>>>>>>>>输入您要私聊的用户名称 exit退出私聊模式")
		fmt.Scanln(&toUser)
	}
}

//更新用户名方法
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>请输入用户名")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("更改用户名失败：", err)
		return false
	}
	return true
}
