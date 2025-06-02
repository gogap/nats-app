package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Build information
var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

// NATSClient represents a NATS client with GUI bindings
type NATSClient struct {
	conn          *nats.Conn
	js            jetstream.JetStream
	status        binding.String
	messageCount  binding.Int
	subscriptions map[string]*nats.Subscription
	messages      binding.StringList
	allMessages   []string
	filter        string
	// JetStream data
	streams       []jetstream.StreamInfo
	consumers     []ConsumerInfo
	mu            sync.RWMutex
	refreshJSFunc func()
}

// ConsumerInfo holds consumer information for display
type ConsumerInfo struct {
	Name       string
	StreamName string
	Config     jetstream.ConsumerConfig
}

// NewNATSClient creates a new NATS client instance
func NewNATSClient() *NATSClient {
	status := binding.NewString()
	status.Set("Disconnected")

	return &NATSClient{
		status:        status,
		messageCount:  binding.NewInt(),
		subscriptions: make(map[string]*nats.Subscription),
		messages:      binding.NewStringList(),
		allMessages:   make([]string, 0),
	}
}

// Connect establishes connection to NATS server
func (nc *NATSClient) Connect(url, username, password string) error {
	opts := []nats.Option{
		nats.Name("Fyne NATS Client"),
		nats.ReconnectWait(time.Second * 2),
		nats.MaxReconnects(5),
	}

	if username != "" && password != "" {
		opts = append(opts, nats.UserInfo(username, password))
	}

	conn, err := nats.Connect(url, opts...)
	if err != nil {
		nc.status.Set("Connection Failed")
		return err
	}

	nc.conn = conn

	// Initialize JetStream
	js, err := jetstream.New(conn)
	if err != nil {
		log.Printf("JetStream not available: %v", err)
		nc.js = nil
	} else {
		nc.js = js

		// Auto-refresh JetStream info after connection
		go func() {
			time.Sleep(500 * time.Millisecond) // Wait for UI to be ready
			nc.RefreshJetStreamInfo()

			// Trigger UI refresh if available
			nc.mu.RLock()
			refreshFunc := nc.refreshJSFunc
			nc.mu.RUnlock()

			if refreshFunc != nil {
				refreshFunc()
			}
		}()
	}

	nc.status.Set("Connected")
	return nil
}

// Disconnect closes the NATS connection
func (nc *NATSClient) Disconnect() {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	if nc.conn != nil {
		// Unsubscribe all active subscriptions
		for _, sub := range nc.subscriptions {
			sub.Unsubscribe()
		}
		nc.subscriptions = make(map[string]*nats.Subscription)

		nc.conn.Close()
		nc.conn = nil
	}
	nc.status.Set("Disconnected")
}

// Publish sends a message to the specified subject
func (nc *NATSClient) Publish(subject, message string) error {
	if nc.conn == nil {
		return fmt.Errorf("not connected")
	}
	return nc.conn.Publish(subject, []byte(message))
}

// Subscribe subscribes to messages on the specified subject
func (nc *NATSClient) Subscribe(subject string) error {
	if nc.conn == nil {
		return fmt.Errorf("not connected")
	}

	nc.mu.Lock()
	defer nc.mu.Unlock()

	// Check if already subscribed
	if _, exists := nc.subscriptions[subject]; exists {
		return fmt.Errorf("already subscribed to %s", subject)
	}

	sub, err := nc.conn.Subscribe(subject, func(msg *nats.Msg) {
		timestamp := time.Now().Format("15:04:05")
		formattedMsg := fmt.Sprintf("[%s] %s: %s", timestamp, msg.Subject, string(msg.Data))

		nc.mu.Lock()
		// Add to all messages
		nc.allMessages = append(nc.allMessages, formattedMsg)

		// Limit to 100 messages
		if len(nc.allMessages) > 100 {
			nc.allMessages = nc.allMessages[1:]
		}

		// Apply filter and update display
		nc.applyFilterLocked()
		nc.mu.Unlock()
	})

	if err != nil {
		return err
	}

	nc.subscriptions[subject] = sub
	return nil
}

// Unsubscribe removes subscription from the specified subject
func (nc *NATSClient) Unsubscribe(subject string) error {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	if sub, exists := nc.subscriptions[subject]; exists {
		err := sub.Unsubscribe()
		if err != nil {
			return err
		}
		delete(nc.subscriptions, subject)
	}
	return nil
}

// GetSubscriptions returns list of active subscriptions
func (nc *NATSClient) GetSubscriptions() []string {
	nc.mu.RLock()
	defer nc.mu.RUnlock()

	subjects := make([]string, 0, len(nc.subscriptions))
	for subject := range nc.subscriptions {
		subjects = append(subjects, subject)
	}
	return subjects
}

// ClearMessages clears all messages from the display
func (nc *NATSClient) ClearMessages() {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	nc.allMessages = make([]string, 0)
	nc.messages.Set([]string{})
	nc.messageCount.Set(0)
}

// SetFilter sets the message filter
func (nc *NATSClient) SetFilter(filter string) {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	nc.filter = filter
	nc.applyFilterLocked()
}

// applyFilterLocked applies the current filter to messages (must be called with lock held)
func (nc *NATSClient) applyFilterLocked() {
	var filteredMessages []string

	if nc.filter == "" {
		filteredMessages = nc.allMessages
	} else {
		for _, msg := range nc.allMessages {
			if strings.Contains(strings.ToLower(msg), strings.ToLower(nc.filter)) {
				filteredMessages = append(filteredMessages, msg)
			}
		}
	}

	nc.messages.Set(filteredMessages)
	nc.messageCount.Set(len(filteredMessages))
}

// RefreshJetStreamInfo refreshes the streams and consumers information
func (nc *NATSClient) RefreshJetStreamInfo() error {
	if nc.js == nil {
		return fmt.Errorf("JetStream not available")
	}

	nc.mu.Lock()
	defer nc.mu.Unlock()

	// Clear existing data
	nc.streams = nil
	nc.consumers = nil

	// Get streams
	streamsCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	streamNames := nc.js.ListStreams(streamsCtx)
	for streamInfo := range streamNames.Info() {
		nc.streams = append(nc.streams, *streamInfo)

		// Get consumers for this stream
		stream, err := nc.js.Stream(context.Background(), streamInfo.Config.Name)
		if err != nil {
			log.Printf("Failed to get stream %s: %v", streamInfo.Config.Name, err)
			continue
		}

		consumerNames := stream.ListConsumers(context.Background())
		for consumerInfo := range consumerNames.Info() {
			nc.consumers = append(nc.consumers, ConsumerInfo{
				Name:       consumerInfo.Name,
				StreamName: streamInfo.Config.Name,
				Config:     consumerInfo.Config,
			})
		}
	}

	return nil
}

// GetStreams returns current streams information
func (nc *NATSClient) GetStreams() []jetstream.StreamInfo {
	nc.mu.RLock()
	defer nc.mu.RUnlock()
	return nc.streams
}

// GetConsumers returns current consumers information
func (nc *NATSClient) GetConsumers() []ConsumerInfo {
	nc.mu.RLock()
	defer nc.mu.RUnlock()
	return nc.consumers
}

func main() {
	myApp := app.New()
	myApp.SetIcon(theme.ComputerIcon())
	myWindow := myApp.NewWindow(fmt.Sprintf("NATS Client v%s", Version))
	myWindow.Resize(fyne.NewSize(1000, 700))
	myWindow.CenterOnScreen()

	client := NewNATSClient()

	// Create UI components
	content := createMainUI(client, myWindow)
	myWindow.SetContent(content)

	// Handle window close
	myWindow.SetCloseIntercept(func() {
		client.Disconnect()
		myApp.Quit()
	})

	myWindow.ShowAndRun()
}

func createMainUI(client *NATSClient, window fyne.Window) *fyne.Container {
	// Menu bar
	mainMenu := createMainMenu(window)
	window.SetMainMenu(mainMenu)

	// Connection area - horizontal layout at top
	connectionArea := createConnectionArea(client, window)

	// Create tabs for Publish, Subscribe, and JetStream
	pubSubTabs := container.NewAppTabs(
		container.NewTabItem("Publish", createPublishTabWithOutput(client, window)),
		container.NewTabItem("Subscribe", createSubscribeTabWithOutput(client)),
		container.NewTabItem("JetStream", createJetStreamTab(client, window)),
	)
	pubSubTabs.SetTabLocation(container.TabLocationTop)

	// Status bar
	statusBar := createStatusBar(client)

	// Main layout: Connection at top, tabs below, status at bottom
	return container.NewBorder(
		container.NewVBox(connectionArea, widget.NewSeparator()), // Top
		statusBar, // Bottom
		nil, nil,  // Left, Right
		pubSubTabs, // Center
	)
}

func createMainMenu(window fyne.Window) *fyne.MainMenu {
	// Help menu
	aboutItem := fyne.NewMenuItem("About", func() {
		content := fmt.Sprintf("NATS Client\n\nVersion: %s\nBuild Time: %s\nGo Version: %s\n\nA visual NATS client built with Fyne.",
			Version, BuildTime, GoVersion)
		dialog.ShowInformation("About", content, window)
	})

	helpMenu := fyne.NewMenu("Help", aboutItem)
	return fyne.NewMainMenu(helpMenu)
}

func createConnectionArea(client *NATSClient, window fyne.Window) *fyne.Container {
	// URL entry with authentication support
	urlEntry := widget.NewEntry()
	urlEntry.SetText("nats://localhost:4222")
	urlEntry.SetPlaceHolder("NATS Server URL (supports nats://user:pass@host:port)")

	// Connection status indicator
	statusLabel := widget.NewLabel("")
	statusLabel.Bind(client.status)

	connectBtn := widget.NewButton("Connect", func() {
		err := client.Connect(urlEntry.Text, "", "") // No separate user/pass
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation("Success", "Connected to NATS server", window)
		}
	})
	connectBtn.Importance = widget.HighImportance

	disconnectBtn := widget.NewButton("Disconnect", func() {
		client.Disconnect()
		dialog.ShowInformation("Info", "Disconnected from NATS server", window)
	})

	// Horizontal layout for connection
	return container.NewBorder(
		nil, nil,
		widget.NewLabel("Server:"),
		container.NewHBox(statusLabel, connectBtn, disconnectBtn),
		urlEntry,
	)
}

func createPublishTabWithOutput(client *NATSClient, window fyne.Window) *fyne.Container {
	// Publish controls area
	publishControls := createPublishControls(client, window)

	// Publish output area (for request-reply responses)
	publishOutput := createPublishOutputArea(client)

	// Add padding around content for better spacing
	leftPanel := container.NewPadded(publishControls)
	rightPanel := container.NewPadded(publishOutput)

	// Split horizontally: controls on left, output on right (50/50)
	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.5) // Equal split: 50% each
	return container.NewBorder(nil, nil, nil, nil, split)
}

func createPublishControls(client *NATSClient, window fyne.Window) *fyne.Container {
	// === Message Configuration Group ===
	subjectEntry := widget.NewEntry()
	subjectEntry.SetPlaceHolder("Subject (e.g., test.subject)")

	// Request timeout entry for request-reply
	timeoutEntry := widget.NewEntry()
	timeoutEntry.SetText("5s")
	timeoutEntry.SetPlaceHolder("5s")

	// Mode selection with timeout
	modeSelect := widget.NewSelect(
		[]string{"Publish", "Request-Reply"},
		func(selected string) {
			// Enable/disable timeout field based on mode
			if selected == "Publish" {
				timeoutEntry.Disable()
			} else {
				timeoutEntry.Enable()
			}
		},
	)
	modeSelect.SetSelected("Publish")
	timeoutEntry.Disable()

	// Three parallel rows: subject, mode, timeout
	subjectRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Subject:"),
		nil,
		subjectEntry,
	)

	modeRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Mode:"),
		nil,
		modeSelect,
	)

	timeoutRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Timeout:"),
		nil,
		timeoutEntry,
	)

	configSection := container.NewVBox(
		subjectRow,
		modeRow,
		timeoutRow,
	)

	// === Message Content Group (no title, with scroll) ===
	messageEntry := widget.NewMultiLineEntry()
	messageEntry.SetPlaceHolder("Message content...")
	messageEntry.Wrapping = fyne.TextWrapWord

	// Use scroll container for message entry
	messageScroll := container.NewScroll(messageEntry)
	messageScroll.SetMinSize(fyne.NewSize(0, 200)) // Minimum height

	// === Action Buttons Group ===
	formatBtn := widget.NewButton("Format JSON", func() {
		var jsonData interface{}
		if err := json.Unmarshal([]byte(messageEntry.Text), &jsonData); err != nil {
			dialog.ShowError(fmt.Errorf("invalid JSON: %v", err), window)
		} else {
			formatted, _ := json.MarshalIndent(jsonData, "", "  ")
			messageEntry.SetText(string(formatted))
		}
	})

	clearBtn := widget.NewButton("Clear", func() {
		messageEntry.SetText("")
	})

	sendBtn := widget.NewButton("Send", func() {
		if subjectEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("subject cannot be empty"), window)
			return
		}

		if modeSelect.Selected == "Request-Reply" {
			// TODO: Implement request-reply functionality
			dialog.ShowInformation("Info", "Request-reply mode coming soon!", window)
		} else {
			err := client.Publish(subjectEntry.Text, messageEntry.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("publish failed: %v", err), window)
			} else {
				dialog.ShowInformation("Success", fmt.Sprintf("Published to %s", subjectEntry.Text), window)
			}
		}
	})
	sendBtn.Importance = widget.HighImportance

	buttonSection := container.NewGridWithColumns(3, formatBtn, clearBtn, sendBtn)

	// Main layout with buttons pinned to bottom
	return container.NewBorder(
		container.NewVBox(
			configSection,
			widget.NewSeparator(),
		), // Top
		buttonSection, // Bottom (pinned)
		nil, nil,      // Left, Right
		messageScroll, // Center (expandable)
	)
}

func createPublishOutputArea(client *NATSClient) *fyne.Container {
	// Output area for request-reply responses
	outputList := widget.NewList(
		func() int { return 0 }, // TODO: Implement publish output messages
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrapWord
			return label
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			// TODO: Display request-reply responses
		},
	)

	clearOutputBtn := widget.NewButton("Clear Output", func() {
		// TODO: Clear publish output
	})

	// Response display card
	responseCard := container.NewBorder(
		container.NewBorder(nil, nil, nil, clearOutputBtn),
		nil, nil, nil,
		container.NewScroll(outputList),
	)

	return responseCard
}

func createSubscribeTabWithOutput(client *NATSClient) *fyne.Container {
	// Subscribe controls area
	subscribeControls := createSubscribeControls(client)

	// Subscribe output area (for received messages)
	subscribeOutput := createSubscribeOutputArea(client)

	// Add padding around content for better spacing
	leftPanel := container.NewPadded(subscribeControls)
	rightPanel := container.NewPadded(subscribeOutput)

	// Split horizontally: controls on left, output on right (50/50)
	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.5) // Equal split: 50% each

	return container.NewBorder(nil, nil, nil, nil, split)
}

func createSubscribeControls(client *NATSClient) *fyne.Container {
	// === Subscription Pattern Group ===
	subjectEntry := widget.NewEntry()
	subjectEntry.SetPlaceHolder("Subject to subscribe (e.g., test.*)")

	examples := widget.NewSelect(
		[]string{"test.*", "events.>", "logs.error.*", "metrics.cpu"},
		func(selected string) {
			subjectEntry.SetText(selected)
		},
	)

	// === Active Subscriptions Group (declare list first) ===
	var subscriptionList *widget.List
	subscriptionList = widget.NewList(
		func() int {
			return len(client.GetSubscriptions())
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel(""),
				widget.NewButton("Unsubscribe", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			subjects := client.GetSubscriptions()
			if id < len(subjects) {
				subject := subjects[id]
				container := obj.(*fyne.Container)
				label := container.Objects[1].(*widget.Label)
				button := container.Objects[2].(*widget.Button)

				label.SetText(subject)
				button.OnTapped = func() {
					err := client.Unsubscribe(subject)
					if err != nil {
						log.Printf("Unsubscribe failed: %v", err)
					} else {
						log.Printf("Unsubscribed from %s", subject)
						subscriptionList.Refresh()
					}
				}
			}
		},
	)

	subscribeBtn := widget.NewButton("Subscribe", func() {
		if subjectEntry.Text == "" {
			log.Println("Subject cannot be empty")
			return
		}

		err := client.Subscribe(subjectEntry.Text)
		if err != nil {
			log.Printf("Subscribe failed: %v", err)
		} else {
			log.Printf("Subscribed to %s", subjectEntry.Text)
			subjectEntry.SetText("")
			subscriptionList.Refresh()
		}
	})
	subscribeBtn.Importance = widget.HighImportance

	// Three parallel rows: pattern, examples, subscribe button (similar to publish layout)
	patternRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Pattern:"),
		nil,
		subjectEntry,
	)

	exampleRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Examples:"),
		nil,
		examples,
	)

	patternSection := container.NewVBox(
		patternRow,
		exampleRow,
		subscribeBtn,
	)

	unsubscribeAllBtn := widget.NewButton("Unsubscribe All", func() {
		subjects := client.GetSubscriptions()
		for _, subject := range subjects {
			client.Unsubscribe(subject)
		}
		subscriptionList.Refresh()
	})

	// Use scroll for subscriptions list with proper height
	subscriptionScroll := container.NewScroll(subscriptionList)
	subscriptionScroll.SetMinSize(fyne.NewSize(0, 200))

	subscriptionSection := container.NewVBox(
		widget.NewLabel("Active Subscriptions:"),
		subscriptionScroll,
	)

	// Main layout with unsubscribe button pinned to bottom
	return container.NewBorder(
		container.NewVBox(
			patternSection,
			widget.NewSeparator(),
		), // Top
		unsubscribeAllBtn, // Bottom (pinned)
		nil, nil,          // Left, Right
		subscriptionSection, // Center (expandable)
	)
}

func createSubscribeOutputArea(client *NATSClient) *fyne.Container {
	// Message list for subscription output
	messageList := widget.NewList(
		func() int {
			msgs, _ := client.messages.Get()
			return len(msgs)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrapWord
			return label
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			msgs, _ := client.messages.Get()
			if id < len(msgs) {
				label := obj.(*widget.Label)
				label.SetText(msgs[id])
			}
		},
	)

	// Bind the list to refresh when messages change
	client.messages.AddListener(binding.NewDataListener(func() {
		messageList.Refresh()
		messageList.ScrollToBottom()
	}))

	// === Filter and Controls Group ===
	filterEntry := widget.NewEntry()
	filterEntry.SetPlaceHolder("Filter messages...")
	filterEntry.OnChanged = func(text string) {
		client.SetFilter(text)
	}

	messageCountLabel := widget.NewLabel("")
	messageCountLabel.Bind(binding.IntToStringWithFormat(client.messageCount, "Messages: %d"))

	autoScrollCheck := widget.NewCheck("Auto-scroll", func(checked bool) {
		// TODO: Implement auto-scroll toggle
	})
	autoScrollCheck.SetChecked(true)

	// Fix filter width by using proper layout
	filterSection := container.NewVBox(
		container.NewBorder(
			nil, nil,
			widget.NewLabel("Filter:"),
			container.NewHBox(messageCountLabel, autoScrollCheck),
			filterEntry, // This will take the remaining space
		),
	)

	// === Action Buttons Group (without title) ===
	clearBtn := widget.NewButton("Clear", func() {
		client.ClearMessages()
	})

	pauseBtn := widget.NewButton("Pause", func() {
		// TODO: Implement pause functionality
	})

	exportBtn := widget.NewButton("Export", func() {
		// TODO: Implement export functionality
	})

	// No title for actions as user suggested
	actionSection := container.NewGridWithColumns(3, pauseBtn, exportBtn, clearBtn)

	// === Message Display with proper scroll ===
	messageScroll := container.NewScroll(messageList)
	messageScroll.SetMinSize(fyne.NewSize(0, 300))

	messageSection := container.NewVBox(
		widget.NewLabel("Received Messages:"),
		messageScroll,
	)

	// Main layout with proper sections
	return container.NewBorder(
		container.NewVBox(
			filterSection,
			widget.NewSeparator(),
			actionSection,
			widget.NewSeparator(),
		), // Top
		nil,      // Bottom
		nil, nil, // Left, Right
		messageSection, // Center (expandable)
	)
}

func createStatusBar(client *NATSClient) *fyne.Container {
	statusLabel := widget.NewLabel("")
	statusLabel.Bind(client.status)

	messageCountLabel := widget.NewLabel("")
	messageCountLabel.Bind(binding.IntToString(client.messageCount))

	timeLabel := widget.NewLabel("")
	go func() {
		for {
			timeLabel.SetText(time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(time.Second)
		}
	}()

	return container.NewBorder(
		nil, nil,
		container.NewHBox(
			widget.NewIcon(theme.InfoIcon()),
			statusLabel,
		),
		timeLabel,
		container.NewHBox(
			widget.NewLabel("Messages:"),
			messageCountLabel,
		),
	)
}

func createJetStreamTab(client *NATSClient, window fyne.Window) *fyne.Container {
	// JetStream controls area
	jetStreamControls := createJetStreamControls(client, window)

	// JetStream monitoring area
	jetStreamOutput := createJetStreamOutput(client)

	// Add padding around content for better spacing
	leftPanel := container.NewPadded(jetStreamControls)
	rightPanel := container.NewPadded(jetStreamOutput)

	// Split horizontally: controls on left, output on right (50/50)
	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.5) // Equal split: 50% each
	return container.NewBorder(nil, nil, nil, nil, split)
}

func createJetStreamControls(client *NATSClient, window fyne.Window) *fyne.Container {
	// === Stream Management ===
	streamNameEntry := widget.NewEntry()
	streamNameEntry.SetPlaceHolder("Stream name (e.g., ORDERS)")

	streamSubjectsEntry := widget.NewEntry()
	streamSubjectsEntry.SetPlaceHolder("Subjects (e.g., orders.*, users.created)")

	retentionSelect := widget.NewSelect(
		[]string{"Limits", "Interest", "WorkQueue"},
		nil,
	)
	retentionSelect.SetSelected("Limits")

	// Stream configuration rows
	streamNameRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Name:"),
		nil,
		streamNameEntry,
	)

	streamSubjectsRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Subjects:"),
		nil,
		streamSubjectsEntry,
	)

	retentionRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Retention:"),
		nil,
		retentionSelect,
	)

	createStreamBtn := widget.NewButton("Create Stream", func() {
		if client.js == nil {
			dialog.ShowError(fmt.Errorf("JetStream not available"), window)
			return
		}

		if streamNameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("stream name cannot be empty"), window)
			return
		}

		subjects := strings.Split(streamSubjectsEntry.Text, ",")
		for i, subject := range subjects {
			subjects[i] = strings.TrimSpace(subject)
		}

		cfg := jetstream.StreamConfig{
			Name:     streamNameEntry.Text,
			Subjects: subjects,
		}

		switch retentionSelect.Selected {
		case "Interest":
			cfg.Retention = jetstream.InterestPolicy
		case "WorkQueue":
			cfg.Retention = jetstream.WorkQueuePolicy
		default:
			cfg.Retention = jetstream.LimitsPolicy
		}

		_, err := client.js.CreateStream(context.Background(), cfg)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to create stream: %v", err), window)
		} else {
			dialog.ShowInformation("Success", fmt.Sprintf("Stream %s created", streamNameEntry.Text), window)
			streamNameEntry.SetText("")
			streamSubjectsEntry.SetText("")

			// Auto-refresh after creating stream
			go func() {
				time.Sleep(100 * time.Millisecond)
				client.mu.RLock()
				refreshFunc := client.refreshJSFunc
				client.mu.RUnlock()

				if refreshFunc != nil {
					refreshFunc()
				}
			}()
		}
	})
	createStreamBtn.Importance = widget.HighImportance

	streamSection := container.NewVBox(
		widget.NewLabel("Stream Management:"),
		streamNameRow,
		streamSubjectsRow,
		retentionRow,
		createStreamBtn,
	)

	// === Consumer Management ===
	consumerNameEntry := widget.NewEntry()
	consumerNameEntry.SetPlaceHolder("Consumer name (e.g., processor)")

	consumerStreamEntry := widget.NewEntry()
	consumerStreamEntry.SetPlaceHolder("Stream name")

	consumerSubjectEntry := widget.NewEntry()
	consumerSubjectEntry.SetPlaceHolder("Filter subject (optional)")

	consumerNameRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Name:"),
		nil,
		consumerNameEntry,
	)

	consumerStreamRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Stream:"),
		nil,
		consumerStreamEntry,
	)

	consumerSubjectRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Filter:"),
		nil,
		consumerSubjectEntry,
	)

	createConsumerBtn := widget.NewButton("Create Consumer", func() {
		if client.js == nil {
			dialog.ShowError(fmt.Errorf("JetStream not available"), window)
			return
		}

		if consumerNameEntry.Text == "" || consumerStreamEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("consumer name and stream name cannot be empty"), window)
			return
		}

		stream, err := client.js.Stream(context.Background(), consumerStreamEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("stream not found: %v", err), window)
			return
		}

		cfg := jetstream.ConsumerConfig{
			Name: consumerNameEntry.Text,
		}

		if consumerSubjectEntry.Text != "" {
			cfg.FilterSubject = consumerSubjectEntry.Text
		}

		_, err = stream.CreateConsumer(context.Background(), cfg)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to create consumer: %v", err), window)
		} else {
			dialog.ShowInformation("Success", fmt.Sprintf("Consumer %s created", consumerNameEntry.Text), window)
			consumerNameEntry.SetText("")
			consumerStreamEntry.SetText("")
			consumerSubjectEntry.SetText("")

			// Auto-refresh after creating consumer
			go func() {
				time.Sleep(100 * time.Millisecond)
				client.mu.RLock()
				refreshFunc := client.refreshJSFunc
				client.mu.RUnlock()

				if refreshFunc != nil {
					refreshFunc()
				}
			}()
		}
	})
	createConsumerBtn.Importance = widget.HighImportance

	consumerSection := container.NewVBox(
		widget.NewLabel("Consumer Management:"),
		consumerNameRow,
		consumerStreamRow,
		consumerSubjectRow,
		createConsumerBtn,
	)

	// === Action Buttons ===
	refreshBtn := widget.NewButton("Refresh", func() {
		client.mu.RLock()
		refreshFunc := client.refreshJSFunc
		client.mu.RUnlock()

		if refreshFunc != nil {
			refreshFunc()
		}
	})

	deleteStreamBtn := widget.NewButton("Delete Stream", func() {
		// TODO: Implement stream deletion with confirmation
	})

	actionSection := container.NewGridWithColumns(2, refreshBtn, deleteStreamBtn)

	// Main layout
	return container.NewBorder(
		container.NewVBox(
			streamSection,
			widget.NewSeparator(),
			consumerSection,
			widget.NewSeparator(),
		), // Top
		actionSection, // Bottom (pinned)
		nil, nil,      // Left, Right
		widget.NewLabel(""), // Center (placeholder)
	)
}

func createJetStreamOutput(client *NATSClient) *fyne.Container {
	// Streams list
	var streamsList *widget.List
	streamsList = widget.NewList(
		func() int {
			return len(client.GetStreams())
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.FolderIcon()),
				widget.NewLabel(""),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			streams := client.GetStreams()
			if id < len(streams) {
				stream := streams[id]
				container := obj.(*fyne.Container)
				nameLabel := container.Objects[1].(*widget.Label)
				infoLabel := container.Objects[2].(*widget.Label)

				nameLabel.SetText(stream.Config.Name)
				infoLabel.SetText(fmt.Sprintf("Msgs: %d, Bytes: %s",
					stream.State.Msgs,
					formatBytes(stream.State.Bytes)))
			}
		},
	)

	streamsScroll := container.NewScroll(streamsList)
	streamsScroll.SetMinSize(fyne.NewSize(0, 200))

	streamsSection := container.NewVBox(
		widget.NewLabel("Streams:"),
		streamsScroll,
	)

	// Consumers list
	var consumersList *widget.List
	consumersList = widget.NewList(
		func() int {
			return len(client.GetConsumers())
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel(""),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			consumers := client.GetConsumers()
			if id < len(consumers) {
				consumer := consumers[id]
				container := obj.(*fyne.Container)
				nameLabel := container.Objects[1].(*widget.Label)
				infoLabel := container.Objects[2].(*widget.Label)

				nameLabel.SetText(consumer.Name)
				infoLabel.SetText(fmt.Sprintf("Stream: %s", consumer.StreamName))
			}
		},
	)

	consumersScroll := container.NewScroll(consumersList)
	consumersScroll.SetMinSize(fyne.NewSize(0, 200))

	consumersSection := container.NewVBox(
		widget.NewLabel("Consumers:"),
		consumersScroll,
	)

	// JetStream info
	jsInfoEntry := widget.NewMultiLineEntry()
	jsInfoEntry.SetPlaceHolder("JetStream information will appear here...")
	jsInfoEntry.Wrapping = fyne.TextWrapWord

	// Update info when refreshed
	updateJSInfo := func() {
		if client.js == nil {
			jsInfoEntry.SetText("JetStream not available")
			return
		}

		info := fmt.Sprintf("JetStream Status: Connected\n")
		info += fmt.Sprintf("Streams: %d\n", len(client.GetStreams()))
		info += fmt.Sprintf("Consumers: %d\n", len(client.GetConsumers()))

		// Add stream details
		for _, stream := range client.GetStreams() {
			info += fmt.Sprintf("\nStream: %s\n", stream.Config.Name)
			info += fmt.Sprintf("  Subjects: %v\n", stream.Config.Subjects)
			info += fmt.Sprintf("  Messages: %d\n", stream.State.Msgs)
			info += fmt.Sprintf("  Bytes: %s\n", formatBytes(stream.State.Bytes))
		}

		jsInfoEntry.SetText(info)
	}

	jsInfoScroll := container.NewScroll(jsInfoEntry)
	jsInfoScroll.SetMinSize(fyne.NewSize(0, 150))

	jsInfoSection := container.NewVBox(
		widget.NewLabel("JetStream Info:"),
		jsInfoScroll,
	)

	// Refresh function
	refreshData := func() {
		err := client.RefreshJetStreamInfo()
		if err != nil {
			log.Printf("Failed to refresh JetStream info: %v", err)
			jsInfoEntry.SetText(fmt.Sprintf("Error: %v", err))
		} else {
			streamsList.Refresh()
			consumersList.Refresh()
			updateJSInfo()
		}
	}

	// Auto-refresh on creation
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to ensure UI is ready
		refreshData()
	}()

	// Store refresh function for external use
	client.mu.Lock()
	client.refreshJSFunc = refreshData
	client.mu.Unlock()

	// Main layout
	return container.NewVBox(
		streamsSection,
		widget.NewSeparator(),
		consumersSection,
		widget.NewSeparator(),
		jsInfoSection,
	)
}

// formatBytes formats byte count as human readable string
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
