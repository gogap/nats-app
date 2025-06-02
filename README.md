# NATS Client - 可视化NATS消息客户端

一个现代化的跨平台NATS客户端，提供简洁易用的图形界面，让NATS消息传递变得更加简单。

[![Download](https://img.shields.io/github/downloads/gogap/nats-app/total?style=for-the-badge&logo=github)](https://github.com/gogap/nats-app/releases)
[![Latest Release](https://img.shields.io/github/v/release/gogap/nats-app?style=for-the-badge)](https://github.com/gogap/nats-app/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

## 🚀 立即下载使用

### 下载安装包

前往 [GitHub Releases](https://github.com/gogap/nats-app/releases/latest) 下载适合您系统的版本：

| 操作系统 | 架构 | 下载链接 |
|---------|------|----------|
| **Windows** | 64位 | [nats-client-windows-amd64.exe](https://github.com/gogap/nats-app/releases/latest) |
| **macOS** | Intel芯片 | [nats-client-darwin-amd64](https://github.com/gogap/nats-app/releases/latest) |
| **macOS** | Apple芯片 | [nats-client-darwin-arm64](https://github.com/gogap/nats-app/releases/latest) |
| **Linux** | 64位 | [nats-client-linux-amd64](https://github.com/gogap/nats-app/releases/latest) |

### 运行应用

下载后双击运行即可，无需额外安装依赖。

## ✨ 主要功能

### 🔗 连接管理
- **智能历史记录**: 自动保存连接地址，支持快速选择
- **一键连接**: 支持 `nats://user:pass@host:port` 格式
- **状态监控**: 实时显示连接状态
- **自动重连**: 连接断开时自动重连

### 📤 消息发布
- **历史记录**: 自动保存发布过的主题，快速重用
- **消息编辑**: 大尺寸编辑器，支持多行消息
- **JSON格式化**: 一键美化JSON消息格式
- **请求-响应**: 支持Request-Reply模式（即将推出）

### 📥 消息订阅
- **模式匹配**: 支持通配符订阅（如 `test.*`, `events.>`）
- **分组订阅**: 支持队列组负载均衡
- **实时接收**: 消息实时显示，自动滚动
- **历史管理**: 保存订阅模式和分组历史

### 💾 JetStream流处理
- **流管理**: 创建和管理JetStream数据流
- **消费者管理**: 配置消息消费者
- **保留策略**: 支持限制、兴趣和工作队列策略
- **实时监控**: 查看流和消费者状态

### 📊 消息管理
- **实时过滤**: 输入关键词即时过滤消息
- **消息统计**: 显示接收消息数量
- **历史记录**: 保留最近100条消息
- **一键清空**: 快速清理消息历史

## 🎯 快速开始

### 1. 启动NATS服务器

如果您还没有NATS服务器，可以：

**使用Docker (推荐):**
```bash
docker run -p 4222:4222 nats:latest
```

**或下载NATS服务器:**
- 前往 [NATS官网](https://nats.io/download/) 下载
- 解压后运行 `nats-server`

### 2. 连接到服务器

1. 打开NATS Client应用
2. 在连接栏输入服务器地址（默认：`nats://localhost:4222`）
3. 点击"连接"按钮

### 3. 发送第一条消息

1. 切换到"发布"标签页
2. 在主题栏输入：`test.hello`
3. 在消息框输入：`Hello, NATS!`
4. 点击"发送"按钮

### 4. 接收消息

1. 切换到"订阅"标签页
2. 在模式栏输入：`test.*`
3. 点击"订阅"按钮
4. 回到发布页面发送消息，即可在订阅页面看到接收的消息

## 📖 使用指南

### 连接配置

支持多种连接格式：
- 基本连接：`nats://localhost:4222`
- 带认证：`nats://username:password@server:4222`
- 集群连接：`nats://server1:4222,server2:4222`

### 消息发布技巧

- **JSON消息**: 粘贴JSON后点击"格式化JSON"美化代码
- **消息模板**: 常用消息会自动保存在历史记录中
- **批量发送**: 可以快速选择历史主题重复发送

### 订阅模式示例

| 模式 | 说明 | 示例 |
|------|------|------|
| `test.*` | 匹配单层级 | `test.hello`, `test.world` |
| `events.>` | 匹配多层级 | `events.user.login`, `events.order.created` |
| `logs.error.*` | 特定层级 | `logs.error.app`, `logs.error.db` |

### JetStream流处理

1. **创建流**: 指定流名称和捕获的主题模式
2. **设置保留策略**:
   - **限制策略**: 按大小/时间限制
   - **兴趣策略**: 有消费者时保留
   - **工作队列**: 确认后删除
3. **创建消费者**: 从流中消费消息

## 🛠️ 配置文件

应用会自动在以下位置保存配置：
- **Windows**: `%APPDATA%\nats-app\config.json`
- **macOS**: `~/Library/Application Support/nats-app/config.json`
- **Linux**: `~/.config/nats-app/config.json`

配置包含：
- 连接历史记录（最多15个）
- 发布主题历史
- 订阅模式历史
- 分组名称历史

## 💡 使用技巧

### 提高效率
- 所有输入框都支持历史记录下拉选择
- 使用快捷模式示例快速订阅常见模式
- 利用分组订阅实现负载均衡
- 定期清理消息历史保持界面整洁

### 调试技巧
- 使用 `>` 通配符监听所有消息进行调试
- 利用消息过滤快速找到特定消息
- 通过JetStream监控查看消息流量统计

### 生产环境
- 使用分组订阅实现服务的高可用
- 通过JetStream实现消息持久化
- 利用流的保留策略控制存储成本

## ❓ 常见问题

**Q: 如何连接到远程NATS服务器？**
A: 在连接栏输入完整URL，如 `nats://192.168.1.100:4222`

**Q: 支持TLS安全连接吗？**
A: 目前版本专注于基础功能，TLS支持在开发路线图中

**Q: 消息历史能保存多少条？**
A: 为保证性能，应用保留最近100条消息，可随时清空

**Q: 配置文件在哪里？**
A: 自动保存在系统配置目录，支持跨平台，详见上方配置文件说明

**Q: 如何报告问题或建议功能？**
A: 请在 [GitHub Issues](https://github.com/gogap/nats-app/issues) 提交

## 🌟 为什么选择这个客户端？

- **简单易用**: 无需复杂配置，下载即用
- **跨平台**: Windows、macOS、Linux原生支持
- **现代界面**: 清晰的标签式界面设计
- **功能完整**: 支持基础消息传递和JetStream高级功能
- **智能历史**: 自动保存使用历史，提高工作效率
- **开源免费**: MIT协议，完全开源

---

## 🔧 开发者信息

### 自行编译

如果您想自己编译应用：

#### 环境要求
- Go 1.21 或更高版本
- Git

#### 编译步骤
```bash
# 克隆项目
git clone https://github.com/gogap/nats-app.git
cd nats-app

# 安装依赖
go mod tidy

# 安装 fyne 打包工具
go install fyne.io/tools/cmd/fyne@latest

# 编译当前平台
fyne package --name nats-client

# 运行 (macOS 会生成 .app 包，Windows 生成 .exe，Linux 生成可执行文件)
# macOS: open nats-client.app
# Windows/Linux: ./nats-client
```

#### 跨平台编译
```bash
# Windows (生成 .exe 文件)
fyne package --os windows --name nats-client-windows

# macOS (生成 .app 应用包)
fyne package --os darwin --name nats-client-darwin

# Linux (生成可执行文件)
fyne package --os linux --name nats-client-linux

# 带版本信息的编译
fyne package --os darwin --name nats-client --app-version 1.0.0 --app-build 1
```

#### 使用Makefile
```bash
make build          # 编译当前平台
make build-all      # 编译所有平台
make clean          # 清理编译文件
```

> **注意**: 使用 `fyne package` 命令可以确保生成的应用程序包含所有必要的资源文件和依赖，在目标平台上可以直接运行，无需安装Go环境。

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [NATS.io](https://nats.io/) - 优秀的消息传递系统
- [Fyne](https://fyne.io/) - 出色的Go GUI框架
- Go社区的支持和贡献

---

**如果觉得有用，请给个 ⭐ Star 支持一下！** 