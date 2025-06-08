package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Client struct {
	ServerIP   string
	ServerPort int
	conn       net.Conn
}

func NewClient(serverIP string, serverPort int) *Client {
	return &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerIP, c.ServerPort))
	if err != nil {
		return fmt.Errorf("连接服务器失败: %v", err)
	}
	c.conn = conn
	return nil
}

func (c *Client) Run() {
	// 启动消息接收goroutine
	go c.ReceiveMessage()

	// 处理用户输入
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("请输入消息（输入'exit'退出，'who'查看在线用户，'to|用户名|消息'私聊）：")
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "exit" {
			break
		}

		if _, err := c.conn.Write([]byte(msg + "\n")); err != nil {
			fmt.Println("发送消息失败:", err)
			return
		}
	}

	// 关闭连接
	c.conn.Close()
	fmt.Println("连接已关闭")
}

func (c *Client) ReceiveMessage() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("使用方法: ./client [服务器IP] [端口]")
		return
	}

	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("无效的端口号:", err)
		return
	}

	client := NewClient(os.Args[1], port)
	if err := client.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	// 处理退出信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		client.conn.Close()
		os.Exit(0)
	}()

	client.Run()
}
