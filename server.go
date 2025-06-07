package main

import (
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

func (s *Server) Handler(conn net.Conn) {
	// ...链接当前的业务
	fmt.Println("Client connected:", conn.RemoteAddr().String())
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
