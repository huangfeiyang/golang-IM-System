package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	Server     string
	ServerPort int
	Name       string
	conn       net.Conn
}

var serverIp string
var serverPort int

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		Server: serverIp,
		ServerPort: serverPort,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	} else {
		client.conn = conn
	}
	return client
}

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "IP地址默认127.0.0.1")
	flag.IntVar(&serverPort, "port", 18888, "端口地址默认18888")

}

func main() {
	flag.Parse()
	client := NewClient("127.0.0.1", 18888)
	if client == nil {
		fmt.Println(">>>>>>>>> 连接服务器失败!")
		return
	} else {
		fmt.Println(">>>>>>>> 连接服务器成功!")
		select {}
	}
}