# NATS Client (Fyne GUI)

A visual NATS client built with the Fyne framework, providing a user-friendly graphical interface for managing NATS connections, publishing, and subscribing to messages.

[![CI](https://github.com/gogap/nats-app/actions/workflows/ci.yml/badge.svg)](https://github.com/gogap/nats-app/actions/workflows/ci.yml)
[![Release](https://github.com/gogap/nats-app/actions/workflows/release.yml/badge.svg)](https://github.com/gogap/nats-app/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gogap/nats-app)](https://goreportcard.com/report/github.com/gogap/nats-app)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Screenshots

*Coming soon - we'll add screenshots of the application interface*

### New Tabbed Interface 🎉

The latest version features an improved UI design with:
- **Tabbed Layout**: Separate tabs for Publish and Subscribe operations
- **More Space**: Larger input areas and better proportions (35% left panel)
- **Enhanced UX**: Pattern examples, bulk operations, and better visual hierarchy

See [docs/UI_DESIGN.md](docs/UI_DESIGN.md) for detailed design documentation.

## Features

### 🔗 Connection Management
- NATS server connection configuration
- Username/password authentication
- Real-time connection status display
- One-click connect/disconnect
- Auto-reconnection with backoff

### 📤 Message Publishing (Enhanced Tab Interface)
- **Larger Editor**: 350x200px message input area
- **Subject Management**: Subject input with history dropdown (planned)
- **Smart Controls**: Format JSON, Clear, and Publish buttons
- **Better Organization**: Clear labels and visual separators
- **Non-destructive**: Preserves message content after publishing

### 📥 Message Subscription (Enhanced Tab Interface)
- **Pattern Examples**: Quick-select dropdown for common patterns (`test.*`, `events.>`, etc.)
- **Enhanced List**: Subscription list with icons and better layout
- **Bulk Operations**: Unsubscribe All button for convenience
- **Visual Indicators**: Clear status and organization
- Support for wildcard subscriptions (e.g., `test.*`, `events.>`)
- Real-time message display
- One-click unsubscribe

### 💬 Message Management (Improved Display)
- **Advanced Toolbar**: Filter, message count, pause/export controls
- **Auto-scroll**: Automatic scrolling to latest messages
- **Text Wrapping**: Better readability with word-wrapped text
- **Real-time Filtering**: Live search as you type
- Message history (up to 100 messages)
- Clear message functionality

### 📊 Status Monitoring
- Connection status indicator
- Message count statistics
- Real-time clock display
- About dialog with version info

### JetStream Management
- ✅ **Stream Management**: Create and manage JetStream streams
- ✅ **Consumer Management**: Create and configure consumers for streams
- ✅ **Retention Policies**: Support for Limits, Interest, and WorkQueue policies
- ✅ **Stream Monitoring**: View stream and consumer information
- ✅ **Subject Filtering**: Advanced filtering for consumers

## Installation

### Download Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/gogap/nats-app/releases):

- **Linux**: `nats-client-linux-amd64`
- **Windows**: `nats-client-windows-amd64.exe`
- **macOS Intel**: `nats-client-darwin-amd64`
- **macOS Apple Silicon**: `nats-client-darwin-arm64`

### Build from Source

First, ensure you have Go 1.21 or higher installed, then:

```bash
git clone https://github.com/gogap/nats-app.git
cd nats-app
go mod tidy
go build -o nats-client .
```

### Using Make

```bash
make build          # Build the application
make run            # Run directly
make build-all      # Build for all platforms
make clean          # Clean build artifacts
```

## Quick Start

1. **Start NATS Server** (if running locally):
   ```bash
   # Install NATS server if you haven't already
   # https://docs.nats.io/running-a-nats-service/introduction/installation
   nats-server
   ```

2. **Run NATS Client**:
   ```bash
   ./nats-client
   ```

3. **Connect**: The default connection URL `nats://localhost:4222` should work for local NATS server

4. **Subscribe**: Try subscribing to `test.*` to receive test messages

5. **Publish**: Send a test message to `test.hello` with content `{"message": "Hello, NATS!"}`

## Usage Guide

For detailed usage instructions, see [docs/USAGE.md](docs/USAGE.md).

### Basic Operations

#### Connect to NATS Server
1. Enter the NATS server address in the "Connection" area (default: `nats://localhost:4222`)
2. If authentication is required, fill in username and password
3. Click "Connect" to establish connection

#### Publish Messages
1. Enter the subject name in the "Publish" area
2. Enter message content in the message box
3. Use "Format JSON" button to format JSON messages
4. Click "Publish" to send the message

#### Subscribe to Messages
1. Enter the subject to subscribe to in the "Subscribe" area
2. Supports wildcards like `test.*` or `events.>`
3. Click "Subscribe" to start subscription
4. Subscribed messages will be displayed in real-time in the right message area

### JetStream Operations

#### Creating Streams
1. Navigate to the "JetStream" tab
2. Enter stream name (e.g., `ORDERS`)
3. Specify subjects to capture (e.g., `orders.*, payments.created`)
4. Select retention policy:
   - **Limits**: Messages retained until size/age limits
   - **Interest**: Messages retained while consumers are interested
   - **WorkQueue**: Messages removed after acknowledgment
5. Click "Create Stream"

#### Creating Consumers
1. Enter consumer name (e.g., `processor`)
2. Specify the stream name to consume from
3. Optionally add a filter subject for selective consumption
4. Click "Create Consumer"

## Configuration

See [examples/config.json](examples/config.json) for example configuration including:
- Connection profiles
- Subject templates
- Message templates

## UI Design Features

### 🎨 Modern Interface
- Native Fyne UI components
- Responsive layout design
- Clear visual hierarchy
- Intuitive user experience

### 📱 Cross-Platform
- Runs on Linux, Windows, and macOS
- Native look and feel on each platform
- Consistent functionality across platforms

### 🔧 Extensibility
- Modular code structure
- Easy to add new features
- Theme customization support

## System Requirements

- **Go**: 1.21+ (for building from source)
- **NATS Server**: Any compatible version (for testing)
- **Operating Systems**:
  - Linux (amd64, arm64)
  - Windows (amd64)
  - macOS (amd64, arm64)

## Development

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Setting up Development Environment

```bash
git clone https://github.com/gogap/nats-app.git
cd nats-app
go mod tidy
```

### Building

```bash
# Build for current platform
go build -o nats-client .

# Or use Make
make build

# Build for all platforms
make build-all
```

### Running Tests

```bash
go test ./...

# Or use Make
make test
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run

# Or use Make
make fmt
make lint
```

## Contributing

We welcome contributions! Please see our contributing guidelines:

### How to Contribute

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add some amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines

- Follow Go best practices and idioms
- Write tests for new functionality
- Update documentation as needed
- Ensure CI checks pass
- Use meaningful commit messages

### Reporting Issues

- Use GitHub Issues for bug reports and feature requests
- Provide detailed information about your environment
- Include steps to reproduce for bug reports
- Check existing issues before creating new ones

## Development Roadmap

- [ ] Message export functionality (JSON, CSV)
- [ ] Connection configuration persistence
- [ ] Subject browser and discovery
- [ ] Message templates and snippets
- [ ] Dark theme support
- [ ] Internationalization (i18n) support
- [ ] Plugin system for extensions
- [ ] Message statistics and analytics
- [ ] TLS/SSL connection support
- [ ] NATS Streaming support

## FAQ

**Q: Why another NATS client?**  
A: We wanted a simple, cross-platform GUI client that's easy to use for development and testing, with modern UI and great user experience.

**Q: Does it support NATS Streaming?**  
A: Not yet, but it's on our roadmap. Currently focuses on core NATS functionality.

**Q: Can I save connection configurations?**  
A: This feature is planned for a future release. Currently, you need to enter connection details each time.

**Q: What's the message limit?**  
A: The application displays up to 100 recent messages to maintain performance. You can clear the history anytime.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [NATS.io](https://nats.io/) for the amazing messaging system
- [Fyne](https://fyne.io/) for the excellent Go GUI framework
- The Go community for the fantastic ecosystem

---

**Star ⭐ this repository if you find it useful!** 