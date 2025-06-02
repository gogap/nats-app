# UI Design Documentation

## Layout Overview (Optimized Design with Visual Grouping)

The NATS Client features a streamlined, efficient interface optimized for real-world NATS usage patterns with enhanced visual hierarchy through proper grouping and spacing.

### Main Window Layout

```
┌─────────────────────────────────────────────────────────────────┐
│  Menu Bar                                                       │
├─────────────────────────────────────────────────────────────────┤
│  Connection: [Server: nats://localhost:4222    ] [●Connected] [Connect] [Disconnect] │
├─────────────────────────────────────────────────────────────────┤
│ ┌─Publish──────────┬─Subscribe──────────────────────────────────┐│
│ │ ┌─ Message Configuration ─────────────────────────────────┐ ││
│ │ │ Subject: [test.subject        ] Mode: [Publish    ▼] │ ││
│ │ │ Timeout: [5s                  ]                       │ ││
│ │ └─────────────────────────────────────────────────────────┘ ││
│ │ ═══════════════════════════════════════════════════════════ ││
│ │ ┌─ Message Content ───────────────────────────────────────┐ ││
│ │ │ ┌─────────────────────────────────────────────────────┐ │ ││
│ │ │ │ Message content...                                  │ │ ││
│ │ │ │                                                     │ │ ││
│ │ │ └─────────────────────────────────────────────────────┘ │ ││
│ │ └─────────────────────────────────────────────────────────┘ ││
│ │ ═══════════════════════════════════════════════════════════ ││
│ │ ┌─ Actions ───────────────────────────────────────────────┐ ││
│ │ │ [Format JSON] [Clear] [Send]                            │ ││
│ │ └─────────────────────────────────────────────────────────┘ ││
│ │─────────────────────────────────────────────────────────────││
│ │ ┌─ Request-Reply Responses ──────────────────────────────┐ ││
│ │ │                                               [Clear] │ ││
│ │ │ (Request-reply response messages)                     │ ││
│ │ └─────────────────────────────────────────────────────────┘ ││
│ └─────────────────┴─────────────────────────────────────────────┘│
├─────────────────────────────────────────────────────────────────┤
│  Status Bar: [●Connected] [Messages: 5] [Time: 15:04:08]       │
└─────────────────────────────────────────────────────────────────┘
```

## Key Design Improvements

### 1. **Visual Grouping with Cards/Frames**

**Problem Solved**: UI elements lacked clear visual hierarchy and organization.

**Solution**:
- **Card-based Layout**: Each functional group wrapped in `widget.NewCard()`
- **Clear Titles**: Each card has descriptive title
- **Visual Separation**: Cards provide natural visual boundaries
- **Consistent Spacing**: Separators between major groups

### 2. **Improved Spacing Strategy**

**Spacing Hierarchy**:
```
Large Gap (Separator): Between major functional groups
Medium Gap (Card padding): Around card content  
Small Gap (Widget spacing): Between related controls
```

**Implementation**:
- `widget.NewSeparator()` between cards
- `widget.NewCard()` provides automatic padding
- Logical grouping reduces cognitive load

### 3. **Enhanced Publish Interface**

**Visual Groups**:
- **Message Configuration**: Subject, mode, timeout
- **Message Content**: Multiline editor with fixed height
- **Actions**: Format, Clear, Send buttons
- **Request-Reply Responses**: Dedicated output area

**Benefits**:
- Clear workflow progression (top to bottom)
- Related controls grouped together
- Visual separation between input and output

### 4. **Optimized Subscribe Interface**

**Visual Groups**:
- **Subscription Pattern**: Pattern input with examples
- **Active Subscriptions**: List with management controls
- **Message Controls**: Filter and display options
- **Actions**: Pause, Export, Clear buttons
- **Received Messages**: Main message display area

**Benefits**:
- Logical workflow organization
- Clear separation of concerns
- Easy subscription management

## Visual Design Principles

### 🎨 **Card-Based Organization**
```go
configCard := widget.NewCard("Message Configuration", "", 
    container.NewVBox(configRow, timeoutRow))
```

### 📏 **Spacing Hierarchy**
```
Card Title
├── Grouped Controls (tight spacing)
├── Related Fields (medium spacing)
└── Action Buttons (grouped)

═══ Separator ═══

Next Card...
```

### 🎯 **Progressive Disclosure**
- Complex features (timeout) only shown when needed
- Mode-aware UI elements
- Clear visual hierarchy guides user attention

## Layout Implementation Details

### Publish Tab Structure
```
PublishTab (VSpilt)
├── PublishControls (30%)
│   ├── Message Configuration Card
│   ├── Separator
│   ├── Message Content Card  
│   ├── Separator
│   └── Actions Card
└── PublishOutput (70%)
    └── Request-Reply Responses Card
```

### Subscribe Tab Structure
```
SubscribeTab (VSplit)
├── SubscribeControls (30%)
│   ├── Subscription Pattern Card
│   ├── Separator
│   └── Active Subscriptions Card
└── SubscribeOutput (70%)
    ├── Message Controls Card
    ├── Separator
    ├── Actions Card
    ├── Separator
    └── Received Messages Card
```

## User Experience Benefits

### 🧠 **Cognitive Load Reduction**
- **Clear Grouping**: Related controls visually grouped
- **Consistent Layout**: Similar patterns across tabs
- **Visual Hierarchy**: Cards guide attention flow

### ⚡ **Workflow Efficiency**
- **Logical Progression**: Top-to-bottom workflow
- **Quick Recognition**: Card titles provide context
- **Reduced Scanning**: Grouped controls easier to find

### 🎨 **Visual Polish**
- **Professional Appearance**: Card-based layout
- **Consistent Spacing**: Systematic spacing rules
- **Clear Boundaries**: Visual separation between areas

## Technical Implementation

### Card Usage Pattern
```go
// Create card with title and content
card := widget.NewCard("Title", "", content)

// Layout with separators for spacing
return container.NewVBox(
    card1,
    widget.NewSeparator(),  // Major separator
    card2,
    widget.NewSeparator(),
    card3,
)
```

### Spacing Strategy
- **Major Groups**: `widget.NewSeparator()` 
- **Card Content**: Automatic card padding
- **Related Controls**: `container.NewVBox()` with minimal spacing
- **Button Groups**: `container.NewGridWithColumns()`

## Best Practices for Contributors

### Visual Grouping Guidelines
1. **Use Cards for Major Groups**: Configuration, Content, Actions
2. **Consistent Titles**: Descriptive card titles
3. **Logical Grouping**: Related controls in same card
4. **Systematic Spacing**: Separators between major groups

### Layout Rules
1. **Top-to-Bottom Flow**: Configuration → Content → Actions → Output
2. **Left-to-Right Priority**: Most important controls on left
3. **Visual Balance**: Consistent card sizes and spacing
4. **Progressive Disclosure**: Hide complexity until needed

### Code Organization
```go
// Pattern: create[Feature][Type]
func createPublishControls() *fyne.Container {
    // Group 1: Configuration
    configCard := widget.NewCard("Configuration", "", ...)
    
    // Group 2: Content  
    contentCard := widget.NewCard("Content", "", ...)
    
    // Group 3: Actions
    actionCard := widget.NewCard("Actions", "", ...)
    
    // Layout with separators
    return container.NewVBox(
        configCard,
        widget.NewSeparator(),
        contentCard, 
        widget.NewSeparator(),
        actionCard,
    )
}
```

This enhanced design provides much better visual organization and user experience through strategic use of cards, proper spacing, and logical grouping of related functionality. 