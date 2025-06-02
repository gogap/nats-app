# Usage Guide

## Getting Started

### 1. Installation

Download the latest release from [GitHub Releases](https://github.com/gogap/nats-app/releases) or build from source:

```bash
git clone https://github.com/gogap/nats-app.git
cd nats-app
go build -o nats-client .
```

### 2. First Launch

Run the application:
```bash
./nats-client
```

The application will open with a default NATS server URL of `nats://localhost:4222`.

## Features

### Connection Management

1. **Server URL**: Enter your NATS server address (e.g., `nats://localhost:4222`)
2. **Authentication**: Optional username and password for secured servers
3. **Connection Status**: Real-time indicator showing connection state
4. **Auto-reconnect**: Automatic reconnection with exponential backoff

### Publishing Messages

1. **Subject**: Enter the subject/topic for your message
2. **Message Content**: Multi-line editor for message body
3. **JSON Formatting**: Click "Format JSON" to pretty-print JSON content
4. **Publish**: Send the message to the specified subject

#### Example Subjects:
- `test.message`
- `events.user.login`
- `metrics.cpu.usage`

### Subscribing to Messages

1. **Subject Pattern**: Enter subject or wildcard pattern
2. **Wildcard Support**:
   - `*` matches a single token (e.g., `events.*` matches `events.login` but not `events.user.login`)
   - `>` matches multiple tokens (e.g., `events.>` matches `events.login` and `events.user.login`)
3. **Subscription Management**: View and unsubscribe from active subscriptions
4. **Real-time Updates**: Messages appear instantly in the message area

#### Example Patterns:
- `test.*` - All test messages
- `events.>` - All event messages
- `logs.error.*` - All error logs
- `metrics.cpu` - Specific CPU metrics

### Message Management

1. **Real-time Display**: Messages appear as they arrive
2. **Filtering**: Use the filter box to search messages by content
3. **Message Limit**: Displays up to 100 recent messages
4. **Clear History**: Remove all displayed messages
5. **Timestamps**: Each message shows arrival time

### Status Monitoring

- **Connection Status**: Shows current connection state
- **Message Count**: Number of received messages
- **Real-time Clock**: Current system time

## Keyboard Shortcuts

- `Ctrl+Enter` (in message editor): Publish message
- `Ctrl+L`: Clear message history
- `Ctrl+F`: Focus filter box
- `F1`: Show about dialog

## Configuration

### Message Templates

You can create predefined message templates for common use cases:

```json
{
  "name": "User Login Event",
  "subject": "events.user.login",
  "body": "{\n  \"user_id\": \"{{USER_ID}}\",\n  \"timestamp\": \"{{TIMESTAMP}}\",\n  \"ip_address\": \"{{IP_ADDRESS}}\"\n}"
}
```

### Connection Profiles

Save frequently used connection settings:

```json
{
  "name": "Production NATS",
  "url": "nats://prod.example.com:4222",
  "username": "client",
  "password": "secret"
}
```

## Common Use Cases

### Development Testing

1. Start a local NATS server: `nats-server`
2. Connect to `nats://localhost:4222`
3. Subscribe to `test.>` to see all test messages
4. Publish test messages to various subjects

### Message Broadcasting

1. Connect to your NATS cluster
2. Subscribe to relevant subjects
3. Use JSON formatting for structured data
4. Monitor message flow in real-time

### System Monitoring

1. Subscribe to `metrics.>` for all metrics
2. Filter by specific services or types
3. Monitor error patterns with `logs.error.*`
4. Set up alerts based on message content

## Troubleshooting

### Connection Issues

- **"Connection Failed"**: Check server URL and network connectivity
- **Authentication errors**: Verify username/password
- **Timeout errors**: Check firewall settings and server availability

### Performance

- **High memory usage**: Clear message history regularly
- **Slow UI**: Reduce number of active subscriptions
- **Missing messages**: Check subject patterns and filters

### Common Problems

1. **No messages appearing**: Verify subscription subject matches published subjects
2. **JSON formatting fails**: Ensure valid JSON syntax
3. **Connection drops**: Check network stability and server configuration

## Tips and Best Practices

1. **Use specific subjects**: Avoid overly broad wildcards in production
2. **Monitor memory**: Clear messages periodically for long-running sessions
3. **Test patterns**: Use local NATS server for testing subject patterns
4. **Organize subjects**: Use hierarchical naming (e.g., `service.component.action`)
5. **JSON validation**: Always validate JSON before publishing critical messages

## Advanced Features

### Message Filtering

- Case-insensitive search
- Partial string matching
- Real-time filtering as you type
- Filter applies to entire message content

### Subscription Management

- View all active subscriptions
- One-click unsubscribe
- Duplicate subscription prevention
- Automatic cleanup on disconnect

### Error Handling

- User-friendly error dialogs
- Detailed error messages
- Graceful degradation
- Automatic recovery where possible 