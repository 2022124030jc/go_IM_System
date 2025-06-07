package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// Online map
	OnlineMap map[string]*User // key: user name, value: User object
	mapLock   sync.Mutex       // Mutex for synchronizing access to OnlineMap
	// Message channel
	Message chan string // channel for broadcasting messages
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User), // Initialize OnlineMap
		mapLock:   sync.Mutex{},           // Initialize mutex
		Message:   make(chan string),      // Initialize message channel
	}
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message // Wait for a message from the channel
		s.mapLock.Lock()   // Lock the map for safe access
		for _, user := range s.OnlineMap {
			user.C <- msg // Send the message to each user
		}
		s.mapLock.Unlock() // Unlock the map after sending messages
	}
}

func (s *Server) Broadcast(user *User, msg string) {
	// Broadcast a message to all users
	s.mapLock.Lock() // Lock the map for safe access
	for _, u := range s.OnlineMap {
		if u.Name != user.Name { // Don't send the message back to the sender
			u.C <- msg // Send the message to the user
		}
	}
	s.mapLock.Unlock() // Unlock the map after broadcasting
}

func (s *Server) Handler(conn net.Conn) {
	// ...链接当前的业务
	fmt.Println("Client connected:", conn.RemoteAddr().String())
	// 创建一个User对象
	user := NewUser(conn)
	// 将用户添加到在线用户列表
	s.mapLock.Lock() // Lock the map for safe access
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock() // Unlock the map after adding the user

}

func (s *Server) Start() {
	// 监听端口
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	// 监听成功后，打印日志
	fmt.Println("Server started at " + s.Ip + ":" + strconv.Itoa(s.Port))

	// 关闭连接
	defer listener.Close()

	// 启动监听消息的协程
	go s.ListenMessage()

	for {
		// accept连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// 启动一个协程来处理连接
		go s.Handler(conn)
	}
}
