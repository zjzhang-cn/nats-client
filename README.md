# NATS Client

一个功能完整的NATS客户端演示项目，包含Go语言后端测试和TypeScript/JavaScript前端Web应用。

## 🚀 特性

### Go 客户端 (后端)
- ✅ NATS连接和基础消息传递
- ✅ JetStream流处理
- ✅ Key-Value存储
- ✅ Object Store对象存储
- ✅ 微服务支持
- ✅ 队列订阅
- ✅ 流发布和拉取订阅
- ✅ HID (硬件接口设备) 集成

### Web 客户端 (前端)
- ✨ **动态服务器配置**: 支持多个预设NATS服务器地址和自定义地址
- ✨ **认证管理**: 可配置的Token认证
- ✨ **实时连接**: WebSocket连接到NATS服务器
- ✨ **JetStream支持**: 流消息的发布和订阅
- ✨ **普通消息**: 传统发布/订阅模式
- ✨ **请求/响应**: 支持同步请求-响应模式
- ✨ **现代UI**: 响应式Web界面

## 📁 项目结构

```
nats-client/
├── go.mod                      # Go模块定义
├── go.sum                      # Go依赖校验和
├── *.go                        # Go源文件
├── *_test.go                   # Go测试文件
├── progress_reader.go          # 进度读取工具
├── html/                       # Web前端应用
│   ├── index.html             # 主页面
│   ├── package.json           # npm包配置
│   ├── vite.config.js         # Vite构建配置
│   ├── src/
│   │   └── index.ts           # TypeScript主文件
│   └── build/                 # 构建输出目录
├── docs/                       # 项目文档
└── CHANGELOG.md               # 变更日志
```

## 🛠️ 技术栈

### 后端 (Go)
- **NATS.go**: NATS Go客户端库
- **Zeroconf**: 服务发现
- **HTTP Proxy**: GE Healthcare内部代理库

### 前端 (TypeScript/JavaScript)
- **NATS.ws**: NATS WebSocket客户端
- **Vite**: 现代前端构建工具
- **TypeScript**: 类型安全的JavaScript

## 🚀 快速开始

### 前置要求

- Go 1.23.7+
- Node.js 16+
- 运行中的NATS服务器 (支持WebSocket)

### 安装和运行

#### 1. 克隆项目
```bash
git clone <repository-url>
cd nats-client
```

#### 2. 运行Go测试
```bash
# 安装Go依赖
go mod download

# 运行所有测试
go test -v ./...

# 运行特定测试
go test -v -run TestNATSConnection
go test -v -run TestJetStreamContext
go test -v -run TestKvWatch
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

## 📚 API文档

### Go客户端主要函数

- `TestNATSConnection()`: 测试基础NATS连接
- `TestJetStreamContext()`: 测试JetStream上下文
- `TestQueueSubscribe()`: 测试队列订阅
- `TestStreamPullSub()`: 测试流拉取订阅
- `TestKvWatch()`: 测试Key-Value监听
- `TestObjectPut/Get()`: 测试对象存储

### Web客户端主要函数

- `connectToNats()`: 连接NATS服务器
- `getSelectedServerUrl()`: 获取选择的服务器地址  
- `getAuthToken()`: 获取认证Token
- `handleServerSelectChange()`: 处理服务器选择变化

## 🔧 开发

### 构建Web应用
```bash
cd html
npm run build
```

构建输出位于 `html/build/` 目录。

### 添加新功能

1. **Go后端**: 在相应的`*_test.go`文件中添加测试用例
2. **Web前端**: 修改`html/src/index.ts`添加新功能
3. **UI更新**: 修改`html/index.html`添加新的界面元素

## 📋 版本历史

查看 [CHANGELOG.md](CHANGELOG.md) 了解详细的版本变更信息。

### 最新版本 v1.1.0 (2025-07-11)
- ✨ 动态NATS服务器地址选择
- ✨ 可配置认证Token
- ✨ 改进的用户界面
- 🔧 增强的连接逻辑和错误处理

## 🤝 贡献

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/新功能`)
3. 提交更改 (`git commit -am '添加新功能'`)
4. 推送到分支 (`git push origin feature/新功能`)
5. 创建Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 支持

如有问题或建议，请：
1. 查看 [文档](docs/)
2. 创建 [Issue](../../issues)
3. 联系项目维护者

---

**作者**: Zhenjiang Zhang  
**版本**: 1.1.0  
**更新时间**: 2025-07-11
