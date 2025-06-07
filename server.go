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
	msg = user.Name + ": " + msg + "\n" // Format the message with the user's name
	s.Message <- msg                    // Send the message to the channel for broadcasting
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

	s.Broadcast(user, "已经上线")

	// 接受客户端消息
	go func() {
		for {
			buf := make([]byte, 4096) // Create a buffer to read data
			n, err := conn.Read(buf)  // Read data from the connection
			if n == 0 || err != nil {
				s.mapLock.Lock()               // Lock the map for safe access
				delete(s.OnlineMap, user.Name) // Remove the user from the online map
				s.mapLock.Unlock()             // Unlock the map after removing the user
				s.Broadcast(user, user.Name+" 已经下线")
				return
			}
			msg := string(buf[:n-1]) // Convert bytes to string
			s.Broadcast(user, msg)   // Broadcast the message
		}
	}()
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
