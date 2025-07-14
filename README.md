# NATS Client

一个功能完整的NATS客户端演示项目，包含Go语言后端测试和TypeScript前端Web应用。本项目展示了如何使用NATS消息系统的各种特性，包括JetStream、Key-Value存储、Object Store等。

## 🚀 特性

### Go 客户端 (后端)
- ✅ **NATS连接管理**: 支持代理连接和环境变量配置
- ✅ **JetStream流处理**: 流发布、拉取订阅和队列订阅
- ✅ **Key-Value存储**: KV存储的更新和监听
- ✅ **Object Store**: 对象的上传和下载
- ✅ **微服务支持**: 基于NATS的微服务架构
- ✅ **队列订阅**: 负载均衡的消息处理
- ✅ **错误处理**: 完整的错误处理和重连机制
- ✅ **进度跟踪**: 文件传输进度监控

### Web 客户端 (前端)
- ✨ **动态服务器配置**: 支持多个预设NATS服务器地址和自定义地址
- ✨ **Token认证**: 可配置的身份验证
- ✨ **WebSocket连接**: 基于nats.ws的浏览器连接
- ✨ **JetStream支持**: 流消息的发布和订阅
- ✨ **普通消息**: 传统发布/订阅模式
- ✨ **请求/响应**: 同步请求-响应通信
- ✨ **现代UI**: 响应式Web界面
- ✨ **实时状态**: 连接状态和消息统计

## 📁 项目结构

```
nats-client/
├── go.mod                      	# Go模块定义
├── go.sum                      	# Go依赖校验和
├── nats_connect.go             	# NATS连接工具
├── progress_reader.go          	# 进度读取工具
├── run.sh                      	# 测试运行脚本
├── *_test.go                   	# 各功能测试文件
│   ├── nats_test.go           		# 基础NATS测试
│   ├── stream-pub_test.go     		# 流发布测试
│   ├── stream-pull-sub_test.go 	# 流拉取订阅测试
│   ├── stream-queue-sub_test.go 	# 流队列订阅测试
│   ├── kv-update_test.go      		# KV更新测试
│   ├── kv-watch_test.go       		# KV监听测试
│   ├── object_put_test.go     		# 对象上传测试
│   ├── object_get_test.go     		# 对象下载测试
│   └── micro_test.go          		# 微服务测试
├── html/                       	# Web前端应用
│   ├── index.html             		# 主页面
│   ├── package.json           		# npm包配置
│   ├── vite.config.js         		# Vite构建配置
│   ├── tsconfig.json          		# TypeScript配置
│   ├── src/
│   │   └── index.ts           		# TypeScript主文件
│   └── build/                 		# 构建输出目录
└── .gitignore                 		# Git忽略文件
```

## 🛠️ 技术栈

### 后端 (Go)
- **Go 1.23.0**: 现代Go语言版本
- **NATS.go v1.40.1**: 官方NATS Go客户端
- **golang.org/x/net**: 网络和代理支持
- **HTTP Proxy**: GE Healthcare内部代理库

### 前端 (TypeScript/JavaScript)
- **NATS.ws v1.30.3**: NATS WebSocket客户端
- **Vite**: 现代前端构建工具
- **TypeScript**: 类型安全的JavaScript

## 🚀 快速开始

### 前置要求

- Go 1.23.0+
- Node.js 16+
- 运行中的NATS服务器 (支持WebSocket)

### 环境配置

配置环境变量:
```bash
export NATS_URL="nats://username:password@your-nats-server:4222"
export ALL_PROXY="socks5://proxy-server:port"  # 可选: 代理配置
```

### 安装和运行

#### 1. 克隆项目
```bash
git clone https://github.com/zjzhang-cn/nats-client.git
cd nats-client
```

#### 2. 运行Go测试
```bash
# 安装Go依赖
go mod download

# 使用脚本运行测试 (包含环境配置)
./run.sh

# 运行所有测试
go test -v ./...

# 运行特定测试模式
./run.sh "TestNATS.*"           # NATS基础功能
./run.sh "TestStream.*"         # JetStream流
./run.sh "TestKv.*"             # Key-Value存储
./run.sh "TestObject.*"         # Object Store
./run.sh "TestMicro.*"          # 微服务
```

#### 3. 启动Web应用
```bash
# 进入前端目录
cd html

# 安装依赖
npm install

# 开发模式启动
npm start

# 构建生产版本
npm run build
```

### 访问Web应用

开发模式下访问: `http://localhost:5173`

## 📖 使用指南

### Web客户端使用

1. **选择服务器地址**
   - 从下拉菜单选择预配置的服务器
   - 或选择"自定义地址..."输入完整WebSocket URL

2. **配置认证**
   - 在Token输入框中输入认证信息

3. **连接服务器**
   - 点击"连接到NATS"建立连接
   - 连接成功后可进行消息操作

4. **消息操作**
   - **JetStream**: 发送和接收流消息
   - **普通订阅**: 订阅主题接收消息
   - **发布消息**: 向指定主题发送消息
   - **请求/响应**: 发送请求并等待响应

### Go客户端测试

项目包含多个测试场景：

- `nats_test.go`: 基础NATS连接和消息测试
- `stream-*_test.go`: JetStream流处理测试
- `kv-*_test.go`: Key-Value存储测试
- `object_*_test.go`: Object Store测试
- `micro_test.go`: 微服务测试
- `hid_test.go`: 硬件接口设备测试

## ⚙️ 配置

### 默认NATS服务器地址

Web客户端预配置的服务器地址：
- `ws://10.189.140.106:1081` (默认)
- `ws://10.189.140.106:1082`
- 自定义地址支持

### 认证配置

- 支持运行时修改认证信息

## 🧪 测试

### 运行Go测试
```bash
# 运行所有测试
go test -v ./...

# 运行基准测试
go test -bench=. -v

# 测试覆盖率
go test -cover ./...
```

### Web应用测试
```bash
cd html

# 开发调试模式
npm run debug

# 构建测试
npm run build
```

## 📚 功能详解

### Go测试模块

#### 基础功能测试 (`nats_test.go`)
- **连接测试**: 验证NATS服务器连接
- **JetStream上下文**: 测试JetStream管理功能
- **队列订阅**: 负载均衡的消息处理
- **多消息处理**: 批量消息发送和接收
- **错误处理**: 连接失败和重连机制

#### 流处理测试
- **`stream-pub_test.go`**: JetStream消息发布
- **`stream-pull-sub_test.go`**: 拉取式消息订阅
- **`stream-queue-sub_test.go`**: 队列模式的流订阅

#### 存储测试
- **`kv-update_test.go`**: Key-Value存储更新操作
- **`kv-watch_test.go`**: Key-Value变化监听
- **`object_put_test.go`**: 对象上传功能
- **`object_get_test.go`**: 对象下载功能

#### 微服务测试 (`micro_test.go`)
- 基于NATS的微服务架构演示

### Web客户端功能

#### 核心特性
- **动态服务器配置**: 支持运行时切换NATS服务器
- **WebSocket连接**: 基于nats.ws的浏览器连接
- **实时状态显示**: 连接状态和消息统计
- **多种消息模式**: JetStream、普通订阅、请求响应

#### 用户界面
- **服务器选择**: 下拉菜单选择或自定义输入
- **认证配置**: Token输入和管理
- **消息面板**: 实时显示接收的消息
- **操作按钮**: 连接控制和消息发送

## 🔧 开发指南

### 环境设置

```bash
# 设置Go代理 (中国用户)
go env -w GOPROXY=https://goproxy.cn,direct

# 安装依赖
go mod tidy
```

### 测试开发

```bash
# 运行单个测试
go test -v -run TestNATSConnection

# 运行所有流相关测试
go test -v -run "TestStream.*"

# 运行测试并显示详细输出
go test -v -run TestKvWatch -args -test.v
```

### Web开发

```bash
cd html

# 开发模式 (热重载)
npm start

# 生产构建
npm run build

# 调试模式
npm run debug
```

## 🚀 部署

### 生产环境部署

1. **构建Web应用**
```bash
cd html
npm run build
```

2. **配置服务器**
- 确保NATS服务器支持WebSocket
- 配置适当的CORS策略
- 设置SSL证书 (生产环境推荐)

3. **部署静态文件**
- 将`html/build/`目录内容部署到Web服务器
- 配置适当的HTTP头部

## � 性能优化

### Go客户端优化
- 使用连接池管理多个NATS连接
- 实现消息批处理减少网络开销
- 配置适当的重连策略

### Web客户端优化
- 启用消息压缩
- 实现消息缓冲和批处理
- 使用Service Worker缓存静态资源

## 🔐 安全考虑

- **认证**: 使用强密码和Token
- **加密**: 生产环境使用TLS连接
- **授权**: 配置适当的NATS权限
- **网络**: 在受信任的网络环境中部署

## 📋 故障排除

### 常见问题

1. **连接失败**
   - 检查NATS服务器是否运行
   - 验证网络连接和防火墙设置
   - 确认WebSocket支持已启用

2. **认证错误**
   - 检查Token格式和有效性
   - 验证服务器认证配置

3. **消息丢失**
   - 检查JetStream配置
   - 验证流的持久化设置

## 📞 支持与反馈

如有问题或建议，请：
- 查看项目文档
- 创建GitHub Issue
- 联系项目维护者: zjzhang-cn

---

**项目**: NATS Client Demo  
**作者**: Zhenjiang Zhang  
**仓库**: https://github.com/zjzhang-cn/nats-client  
**最后更新**: 2025-07-14
