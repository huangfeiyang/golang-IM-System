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

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	//消息广播的channel
	Message chan string
}

//创建一个sever的接口

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}

	return server
}

//用户上线消息发送到server的channel
func (this *Server) BroadCast(user *User, msg string)  {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	//...当前连接的业务
	fmt.Println("连接建立成功!")

	//用户上线,map加入用户
	user := NewUser(conn, this)
	user.Online()
	
	//监听用户活跃状态
	isLive := make(chan bool)

	//接受客户端传递发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用户的消息(去除'\n')
			msg := string(buf[:n-1])

			user.Domessage(msg)
			isLive <- true
		}
	}()
	//当前handler阻塞
	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 10):
				user.SendMsg("登陆状态超时，已强制下线")
				close(user.C)
				conn.Close()
				return
		}
	}
}

func (this *Server) Start() {
	//socket listen
	listener, error := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if error != nil {
		fmt.Println("net.listen err:", error)
		return
	}

	//close listen socket
	defer listener.Close()

	go this.ListenMessager()

	// fmt.Println(listener)
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accecpt err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}

}
