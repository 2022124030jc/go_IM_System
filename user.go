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

// 给当前用户发消息
func (this *User) SendMsg(msg string) {
	if _, err := this.conn.Write([]byte(msg)); err != nil {
		// 如果发送失败，可能是用户已经下线
		this.Offline()
	}
}

// 用户处理消息业务
func (this *User) DoMessage(msg string) {
	// 如果消息是"exit"，则下线
	if msg == "exit" {
		this.Offline()
		return
	} else if msg == "who" {
		// 如果消息是"who"，则返回在线用户列表
		this.server.mapLock.Lock()

		for _, user := range this.server.OnlineMap {
			// 将在线用户的名称发送给当前用户
			this.SendMsg("[" + user.Addr + "] :" + user.Name + "已上线\n")
		}

		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" { //更新用户名
		// 如果消息是"rename|新用户名"，则更新用户名
		newName := msg[7:]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			// 如果新用户名已存在，发送错误消息
			this.SendMsg("用户名已存在，请重新输入\n")
			return
		} else {
			// 如果新用户名不存在，更新用户名
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name) // 删除旧用户名
			this.Name = newName                      // 更新用户名
			this.server.OnlineMap[this.Name] = this  // 添加新用户名
			this.server.mapLock.Unlock()

			// 通知其他用户
			this.server.Broadcast(this, "已将用户名改为："+this.Name)
		}
	} else {
		// 否则，将消息广播给其他用户
		this.server.Broadcast(this, msg)
	}
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
