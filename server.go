package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

//创建一个sever的接口

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (thid *Server) Handler(conn net.Conn) {
	//...当前连接的业务
	fmt.Println("连接建立成功!")
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
