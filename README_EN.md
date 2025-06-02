# NATS Client - Visual NATS Messaging Client

**📖 Languages / 语言支持**: [English](README_EN.md) | [中文](README.md)

A modern cross-platform NATS client with a clean and intuitive graphical interface that makes NATS messaging simple and accessible.

[![Download](https://img.shields.io/github/downloads/gogap/nats-app/total?style=for-the-badge&logo=github)](https://github.com/gogap/nats-app/releases)
[![Latest Release](https://img.shields.io/github/v/release/gogap/nats-app?style=for-the-badge)](https://github.com/gogap/nats-app/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

## 🚀 Quick Download & Setup

### Download Binaries

Go to [GitHub Releases](https://github.com/gogap/nats-app/releases/latest) and download the version for your system:

| Operating System | Architecture | Download Link |
|------------------|--------------|---------------|
| **Windows** | 64-bit | [nats-client-windows-amd64.exe](https://github.com/gogap/nats-app/releases/latest) |
| **macOS** | Intel | [nats-client-darwin-amd64](https://github.com/gogap/nats-app/releases/latest) |
| **macOS** | Apple Silicon | [nats-client-darwin-arm64](https://github.com/gogap/nats-app/releases/latest) |
| **Linux** | 64-bit | [nats-client-linux-amd64](https://github.com/gogap/nats-app/releases/latest) |

### Run the Application

Simply double-click the downloaded file to run - no additional dependencies required.

## ✨ Key Features

### 🔗 Connection Management
- **Smart History**: Automatically saves connection URLs with quick selection
- **One-click Connect**: Supports `nats://user:pass@host:port` format
- **Status Monitoring**: Real-time connection status display
- **Auto-reconnect**: Automatic reconnection on disconnection

### 📤 Message Publishing
- **History Records**: Automatically saves published subjects for quick reuse
- **Message Editor**: Large-size editor supporting multi-line messages
- **JSON Formatting**: One-click JSON beautification
- **Request-Reply**: Support for Request-Reply pattern (coming soon)

### 📥 Message Subscription
- **Pattern Matching**: Support wildcard subscriptions (e.g., `test.*`, `events.>`)
- **Group Subscriptions**: Support queue groups for load balancing
- **Real-time Reception**: Messages displayed in real-time with auto-scroll
- **History Management**: Save subscription patterns and group history

### 💾 JetStream Processing
- **Stream Management**: Create and manage JetStream data streams
- **Consumer Management**: Configure message consumers
- **Retention Policies**: Support Limits, Interest, and WorkQueue policies
- **Real-time Monitoring**: View stream and consumer status

### 📊 Message Management
- **Real-time Filtering**: Instant keyword filtering as you type
- **Message Statistics**: Display received message count
- **Message History**: Retains latest 100 messages
- **One-click Clear**: Quick message history cleanup

## 🎯 Quick Start

### 1. Start NATS Server

If you don't have a NATS server yet:

**Using Docker (Recommended):**
```bash
docker run -p 4222:4222 nats:latest
```

**Or download NATS server:**
- Visit [NATS Official Website](https://nats.io/download/) to download
- Extract and run `nats-server`

### 2. Connect to Server

1. Open the NATS Client application
2. Enter server address in connection bar (default: `nats://localhost:4222`)
3. Click "Connect" button

### 3. Send Your First Message

1. Switch to "Publish" tab
2. Enter subject: `test.hello`
3. Enter message: `Hello, NATS!`
4. Click "Send" button

### 4. Receive Messages

1. Switch to "Subscribe" tab
2. Enter pattern: `test.*`
3. Click "Subscribe" button
4. Go back to publish page and send messages - you'll see them in the subscribe page

## 📖 User Guide

### Connection Configuration

Supports multiple connection formats:
- Basic connection: `nats://localhost:4222`
- With authentication: `nats://username:password@server:4222`
- Cluster connection: `nats://server1:4222,server2:4222`

### Message Publishing Tips

- **JSON Messages**: Paste JSON and click "Format JSON" to beautify
- **Message Templates**: Common messages are automatically saved in history
- **Batch Sending**: Quickly select historical subjects for repeated sending

### Subscription Pattern Examples

| Pattern | Description | Examples |
|---------|-------------|----------|
| `test.*` | Single level match | `test.hello`, `test.world` |
| `events.>` | Multi-level match | `events.user.login`, `events.order.created` |
| `logs.error.*` | Specific level | `logs.error.app`, `logs.error.db` |

### JetStream Stream Processing

1. **Create Stream**: Specify stream name and subject patterns to capture
2. **Set Retention Policy**:
   - **Limits Policy**: Limit by size/time
   - **Interest Policy**: Retain while consumers are interested
   - **Work Queue**: Delete after acknowledgment
3. **Create Consumer**: Consume messages from stream

## 🛠️ Configuration Files

The application automatically saves configuration in:
- **Windows**: `%APPDATA%\nats-app\config.json`
- **macOS**: `~/Library/Application Support/nats-app/config.json`
- **Linux**: `~/.config/nats-app/config.json`

Configuration includes:
- Connection history (up to 15 entries)
- Published subject history
- Subscription pattern history
- Group name history

## 💡 Usage Tips

### Improve Efficiency
- All input fields support history dropdown selection
- Use quick pattern examples for common subscription patterns
- Utilize group subscriptions for load balancing
- Regularly clear message history to keep interface clean

### Debugging Tips
- Use `>` wildcard to monitor all messages for debugging
- Use message filtering to quickly find specific messages
- View message flow statistics through JetStream monitoring

### Production Environment
- Use group subscriptions for high availability services
- Implement message persistence through JetStream
- Control storage costs with stream retention policies

## ❓ Frequently Asked Questions

**Q: How to connect to remote NATS server?**
A: Enter the complete URL in connection bar, e.g., `nats://192.168.1.100:4222`

**Q: Does it support TLS secure connections?**
A: Current version focuses on basic functionality, TLS support is on the roadmap

**Q: How many message history entries can be saved?**
A: For performance, the application retains the latest 100 messages, can be cleared anytime

**Q: Where is the configuration file?**
A: Automatically saved in system configuration directory, cross-platform compatible, see configuration section above

**Q: How to report issues or suggest features?**
A: Please submit at [GitHub Issues](https://github.com/gogap/nats-app/issues)

## 🌟 Why Choose This Client?

- **Easy to Use**: Download and run, no complex configuration needed
- **Cross-platform**: Native support for Windows, macOS, Linux
- **Modern Interface**: Clean tabbed interface design
- **Feature Complete**: Supports basic messaging and JetStream advanced features
- **Smart History**: Automatically saves usage history for improved efficiency
- **Open Source & Free**: MIT license, completely open source

---

## 🔧 Developer Information

### Build from Source

If you want to compile the application yourself:

#### Requirements
- Go 1.21 or higher
- Git

#### Build Steps
```bash
# Clone repository
git clone https://github.com/gogap/nats-app.git
cd nats-app

# Install dependencies
go mod tidy

# Build for current platform
make build

# Run
# macOS: open nats-client.app
# Windows/Linux: ./nats-client
```

#### Font Support
This application uses Go embed technology for built-in Chinese font support:
- **Font**: Source Han Sans CN Medium weight (SourceHanSansCN-Medium.otf)
- **Size**: ~8MB
- **Advantage**: No system font installation required, automatic Chinese display support
- **Implementation**: Auto-embedded during build, no additional steps needed

#### Cross-platform Compilation
```bash
# Windows
fyne package --os windows --name nats-client-windows

# macOS
fyne package --os darwin --name nats-client-darwin

# Linux
fyne package --os linux --name nats-client-linux

# With version information
fyne package --os darwin --name nats-client --app-version 1.0.0 --app-build 1
```

#### Using Makefile
```bash
make build          # Build current platform
make build-release  # Build release version
make clean          # Clean build files
```

> **Note**: Using `fyne package` ensures the generated application includes all necessary resource files and dependencies, and can run directly on target platforms without requiring Go environment. Font files are auto-embedded via Go embed for optimal Chinese display.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [NATS.io](https://nats.io/) - Excellent messaging system
- [Fyne](https://fyne.io/) - Outstanding Go GUI framework
- Go community support and contributions

---

**If you find this useful, please give it a ⭐ Star!** 