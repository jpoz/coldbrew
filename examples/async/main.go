package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jpoz/trmnl"
)

// AsyncModel represents our async demo application state
type AsyncModel struct {
	Status      string
	Data        []string
	Loading     bool
	Error       string
	RequestID   int
	LastUpdate  time.Time
}

// Message types
type LoadDataMsg struct{}
type DataLoadedMsg struct {
	RequestID int
	Data      []string
}
type LoadErrorMsg struct {
	RequestID int
	Error     string
}
type TickMsg struct {
	Time time.Time
}

// Init initializes the model
func (m AsyncModel) Init() (trmnl.Model, trmnl.Cmd) {
	m.Status = "Ready"
	m.Data = []string{}
	m.Loading = false
	m.LastUpdate = time.Now()
	return m, nil
}

// Update handles messages and returns updated model
func (m AsyncModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case trmnl.KeyMsg:
		switch msg.String() {
		case "q":
			return m, trmnl.Quit()
		case "l", "load":
			if !m.Loading {
				m.Loading = true
				m.Status = "Loading data..."
				m.Error = ""
				m.RequestID++
				return m, loadDataCmd(m.RequestID)
			}
		case "c", "clear":
			m.Data = []string{}
			m.Status = "Data cleared"
			m.Error = ""
		default:
			m.Status = fmt.Sprintf("Unknown key: %s. Press 'l' to load, 'c' to clear, 'q' to quit", msg.String())
		}
		
	case LoadDataMsg:
		if !m.Loading {
			m.Loading = true
			m.Status = "Loading data..."
			m.Error = ""
			m.RequestID++
			return m, loadDataCmd(m.RequestID)
		}
		
	case DataLoadedMsg:
		// Only update if this is the current request
		if msg.RequestID == m.RequestID {
			m.Loading = false
			m.Data = msg.Data
			m.Status = fmt.Sprintf("Loaded %d items", len(msg.Data))
			m.Error = ""
			m.LastUpdate = time.Now()
		}
		
	case LoadErrorMsg:
		// Only update if this is the current request
		if msg.RequestID == m.RequestID {
			m.Loading = false
			m.Error = msg.Error
			m.Status = "Error loading data"
		}
		
	case TickMsg:
		// Just update the time for display
		m.LastUpdate = msg.Time
	}
	
	return m, nil
}

// View renders the model to a component tree
func (m AsyncModel) View() trmnl.Component {
	// Create title
	title := trmnl.NewText("Async Operations Demo")
	title.Style.Border = true
	title.Style.RoundedBorder = true
	title.Style.BorderColor = trmnl.ColorBlue
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create status display
	statusColor := trmnl.ColorGreen
	if m.Loading {
		statusColor = trmnl.ColorYellow
	} else if m.Error != "" {
		statusColor = trmnl.ColorRed
	}
	
	statusText := m.Status
	if m.Loading {
		statusText += " " + getSpinner()
	}
	
	status := trmnl.NewText(statusText)
	status.Style.Border = true
	status.Style.BorderColor = statusColor
	status.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create data display
	dataText := fmt.Sprintf("Data Items (%d):", len(m.Data))
	if len(m.Data) == 0 {
		dataText += "\n  (no data)"
	} else {
		for i, item := range m.Data {
			if i < 5 { // Show only first 5 items
				dataText += fmt.Sprintf("\n  %d. %s", i+1, item)
			} else if i == 5 {
				dataText += fmt.Sprintf("\n  ... and %d more", len(m.Data)-5)
				break
			}
		}
	}
	
	dataDisplay := trmnl.NewText(dataText)
	dataDisplay.Style.Border = true
	dataDisplay.Style.BorderColor = trmnl.ColorCyan
	dataDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create error display if there's an error
	var errorDisplay *trmnl.Text
	if m.Error != "" {
		errorText := fmt.Sprintf("Error: %s", m.Error)
		errorDisplay = trmnl.NewText(errorText)
		errorDisplay.Style.Border = true
		errorDisplay.Style.BorderColor = trmnl.ColorRed
		errorDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	}
	
	// Create instructions
	instructions := trmnl.NewText("Commands:\n  l : Load data\n  c : Clear data\n  q : Quit")
	instructions.Style.Border = true
	instructions.Style.BorderColor = trmnl.ColorMagenta
	instructions.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create info
	info := trmnl.NewText(fmt.Sprintf("Last Update: %s\nRequest ID: %d", 
		m.LastUpdate.Format("15:04:05"), m.RequestID))
	info.Style.Border = true
	info.Style.BorderColor = trmnl.ColorBrightBlack
	info.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Layout everything
	container := trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignCenter).
		SetPadding(trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}).
		AddChild(title).
		AddChild(status).
		AddChild(dataDisplay)
	
	if errorDisplay != nil {
		container.AddChild(errorDisplay)
	}
	
	container.AddChild(instructions).AddChild(info)
	
	return container
}

// Subscriptions returns the subscriptions for this model
func (m AsyncModel) Subscriptions() []trmnl.Sub {
	return []trmnl.Sub{
		trmnl.Every(time.Second, func(t time.Time) trmnl.Msg {
			return TickMsg{Time: t}
		}),
	}
}

// loadDataCmd simulates loading data asynchronously
func loadDataCmd(requestID int) trmnl.Cmd {
	return func() trmnl.Msg {
		// Simulate network delay
		delay := time.Duration(rand.Intn(3)+1) * time.Second
		time.Sleep(delay)
		
		// Simulate random success/failure
		if rand.Float64() < 0.2 { // 20% chance of failure
			return LoadErrorMsg{
				RequestID: requestID,
				Error:     "Network timeout or server error",
			}
		}
		
		// Generate fake data
		items := []string{
			"User Profile Data",
			"Recent Messages",
			"System Configuration",
			"Analytics Report",
			"Cache Statistics",
			"Performance Metrics",
			"Security Logs",
			"Database Schema",
			"API Endpoints",
			"Feature Flags",
		}
		
		// Return random subset
		count := rand.Intn(8) + 3 // 3-10 items
		selectedItems := make([]string, count)
		for i := 0; i < count; i++ {
			selectedItems[i] = items[rand.Intn(len(items))]
		}
		
		return DataLoadedMsg{
			RequestID: requestID,
			Data:      selectedItems,
		}
	}
}

// getSpinner returns a simple spinner animation
func getSpinner() string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return spinners[int(time.Now().UnixNano()/100000000)%len(spinners)]
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Create initial model
	initialModel := AsyncModel{
		Status:     "Ready to load data",
		Data:       []string{},
		Loading:    false,
		RequestID:  0,
		LastUpdate: time.Now(),
	}
	
	// Create and run program
	program := trmnl.NewProgram(initialModel)
	if err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}