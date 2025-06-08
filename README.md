# Go Instant Messaging System

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-blue)

基于Go语言实现的即时通讯系统，包含服务端和客户端组件，支持多人聊天、私聊、用户状态管理等基础IM功能。

## 功能特性

- **用户管理**
  - 自动识别新连接用户
  - 用户名修改 (`rename|新用户名`)
  - 用户上下线广播通知
  - 在线用户查询 (`who`)

- **消息功能**
  - 公聊消息广播
  - 私聊消息 (`to|用户名|消息内容`)
  - 消息格式校验
  - 离线消息处理

- **系统特性**
  - 并发连接处理
  - 连接状态监控
  - 异常断开自动清理
  - 双通道消息队列

## 技术栈

- **核心语言**: Go 1.21+
- **网络通信**: net 标准库
- **并发模型**: goroutine + channel
- **数据同步**: sync.Mutex
- **终端处理**: bufio.Scanner

## 快速开始

### 服务端运行

```bash
# 进入项目目录
cd go_IM_System

# 启动服务器 (默认端口8888)
go run main.go
```

### 客户端使用

```bash
# 编译客户端
cd client && go build -o im_client .

# 运行客户端
./im_client 127.0.0.1 8888

# 连接成功后可使用以下命令：
# 查看在线用户：who
# 私聊用户：to|用户名|消息内容  
# 修改用户名：rename|新用户名
# 退出系统：exit
```

## 项目结构

```
go_IM_System/
├── client/               # 客户端实现
│   └── client.go         # 客户端主逻辑
├── server.go             # 服务端主程序
├── user.go               # 用户连接管理
├── main.go               # 服务端入口
├── go.mod                # 依赖管理
└── README.md             # 项目文档
```

## 贡献指南

1. 提交Issue描述问题或建议
2. Fork仓库并创建功能分支
3. 提交Pull Request时关联相关Issue
4. 遵循现有代码风格和提交规范

## 许可证

[MIT License](LICENSE) © 2022124030jc
