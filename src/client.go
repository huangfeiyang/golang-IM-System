package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	Server     string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

var serverIp string
var serverPort int

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		Server:     serverIp,
		ServerPort: serverPort,
		flag:       999,
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

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>>>>>> 请输入合法数字！")
		return false
	}
}

func (client * Client) UpdateName() bool {
	fmt.Println("请输入用户名:")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
	//以上等价于以下
	// for {
	// 	buf := make()
	// 	client.conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println("您已进入公聊模式!")
			break
		case 2:
			fmt.Println("您已进入私聊模式")
			client.PublicChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "IP地址默认127.0.0.1")
	flag.IntVar(&serverPort, "port", 18888, "端口地址默认18888")

}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient("127.0.0.1", 18888)
	if client == nil {
		fmt.Println(">>>>>>>>> 连接服务器失败!")
		return
	} else {

		//开启go程处理相应
		go client.DealResponse()

		fmt.Println(">>>>>>>> 连接服务器成功!")
		client.Run()
	}
}
