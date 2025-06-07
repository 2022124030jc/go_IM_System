package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server // 关联Server对象
}

// 创建用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server, // 关联Server对象
	}

	// 启动一个协程来监听当前用户的消息
	go user.ListenMessage()

	return user
}

// 用户上线业务
func (this *User) Online() {
	// 将用户添加到在线用户列表
	this.server.mapLock.Lock() // Lock the map for safe access
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock() // Unlock the map after adding the user

	this.server.Broadcast(this, "已经上线")
}

// 用户下线业务
func (this *User) Offline() {
	// 将用户从在线用户列表中删除
	this.server.mapLock.Lock() // Lock the map for safe access
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock() // Unlock the map after removing the user

	this.server.Broadcast(this, "已经下线")
}

// 用户处理消息业务
func (this *User) DoMessage(msg string) {
	// 如果消息是"exit"，则下线
	if msg == "exit" {
		this.Offline()
		return
	}

	// 否则，将消息广播给其他用户
	this.server.Broadcast(this, msg)
}

// 监听当前User的消息
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		if _, err := this.conn.Write([]byte(msg)); err != nil {
			return
		}
	}
}
