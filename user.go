package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn
}

// 创建用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
	}

	// 启动一个协程来监听当前用户的消息
	go user.ListenMessage()

	return user
}

// 监听当前User的消息
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		if _, err := this.Conn.Write([]byte(msg)); err != nil {
			return
		}
	}
}
