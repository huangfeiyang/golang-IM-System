package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
	server *Server
}

//创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
		server: server,
	}

	//启动监听当前user的go程
	go user.ListenMessage()
	return user
}

//监听当前user的channel，一旦监听到就发送给client
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}

//用户上线功能
func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "已上线")
}

//用户下线功能
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "已下线")
}

//用户接收消息功能
func (this *User) Domessage(msg string) {
	this.server.BroadCast(this, msg)
}