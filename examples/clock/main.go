package main

import (
	"fmt"
	"time"

	"github.com/jpoz/trmnl"
)

// ClockModel represents our clock application state
type ClockModel struct {
	CurrentTime time.Time
	Running     bool
	Message     string
}

// Message types
type TickMsg struct {
	Time time.Time
}

type ToggleMsg struct{}

// Init initializes the model
func (m ClockModel) Init() (trmnl.Model, trmnl.Cmd) {
	m.CurrentTime = time.Now()
	m.Running = true
	m.Message = "Press 't' to toggle, 'q' to quit"
	return m, nil
}

// Update handles messages and returns updated model
func (m ClockModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case trmnl.KeyMsg:
		switch msg.String() {
		case "q":
			return m, trmnl.Quit()
		case "t":
			m.Running = !m.Running
			if m.Running {
				m.Message = "Clock started. Press 't' to toggle, 'q' to quit"
			} else {
				m.Message = "Clock paused. Press 't' to toggle, 'q' to quit"
			}
		default:
			m.Message = fmt.Sprintf("Unknown key: %s. Press 't' to toggle, 'q' to quit", msg.String())
		}
		
	case TickMsg:
		if m.Running {
			m.CurrentTime = msg.Time
		}
	}
	
	return m, nil
}

// View renders the model to a component tree
func (m ClockModel) View() trmnl.Component {
	// Create title
	title := trmnl.NewText("Digital Clock")
	title.Style.Border = true
	title.Style.RoundedBorder = true
	title.Style.BorderColor = trmnl.ColorBlue
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create time display
	timeFormat := m.CurrentTime.Format("15:04:05")
	dateFormat := m.CurrentTime.Format("Monday, January 2, 2006")
	
	timeText := fmt.Sprintf("%s\n%s", timeFormat, dateFormat)
	clockDisplay := trmnl.NewText(timeText)
	clockDisplay.Style.Border = true
	clockDisplay.Style.BorderColor = trmnl.ColorGreen
	clockDisplay.Style.Padding = trmnl.Box{Top: 2, Right: 4, Bottom: 2, Left: 4}
	
	// Create status
	status := "Running"
	statusColor := trmnl.ColorGreen
	if !m.Running {
		status = "Paused"
		statusColor = trmnl.ColorYellow
	}
	
	statusDisplay := trmnl.NewText(fmt.Sprintf("Status: %s", status))
	statusDisplay.Style.Border = true
	statusDisplay.Style.BorderColor = statusColor
	statusDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create message display
	message := trmnl.NewText(m.Message)
	message.Style.Border = true
	message.Style.BorderColor = trmnl.ColorCyan
	message.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Layout everything
	return trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignCenter).
		SetPadding(trmnl.Box{Top: 2, Right: 2, Bottom: 2, Left: 2}).
		AddChild(title).
		AddChild(clockDisplay).
		AddChild(statusDisplay).
		AddChild(message)
}

// Subscriptions returns the subscriptions for this model
func (m ClockModel) Subscriptions() []trmnl.Sub {
	return []trmnl.Sub{
		trmnl.Every(time.Second, func(t time.Time) trmnl.Msg {
			return TickMsg{Time: t}
		}),
	}
}

func main() {
	// Create initial model
	initialModel := ClockModel{
		CurrentTime: time.Now(),
		Running:     true,
	}
	
	// Create and run program
	program := trmnl.NewProgram(initialModel)
	if err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}