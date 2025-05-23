package main

import (
	"fmt"

	"github.com/jpoz/trmnl"
)

// CursorDemoModel represents our cursor demo application state
type CursorDemoModel struct {
	Message       string
	CursorVisible bool
	Count         int
}

// Message types
type ToggleCursorMsg struct{}
type IncrementMsg struct{}

// Init initializes the model
func (m CursorDemoModel) Init() (trmnl.Model, trmnl.Cmd) {
	m.Message = "Press 'c' to toggle cursor, '+' to increment, 'q' to quit"
	m.CursorVisible = false // Start with cursor hidden
	return m, nil
}

// Update handles messages and returns updated model
func (m CursorDemoModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case trmnl.KeyMsg:
		switch msg.String() {
		case "q":
			return m, trmnl.Quit()
		case "c":
			m.CursorVisible = !m.CursorVisible
			if m.CursorVisible {
				m.Message = "Cursor is now VISIBLE. Press 'c' to hide, '+' to increment, 'q' to quit"
			} else {
				m.Message = "Cursor is now HIDDEN. Press 'c' to show, '+' to increment, 'q' to quit"
			}
			return m, toggleCursorCmd(m.CursorVisible)
		case "+", "=":
			m.Count++
			m.Message = fmt.Sprintf("Count: %d. Press 'c' to toggle cursor, 'q' to quit", m.Count)
		default:
			m.Message = fmt.Sprintf("Unknown key: %s. Press 'c' to toggle cursor, '+' to increment, 'q' to quit", msg.String())
		}
		
	case ToggleCursorMsg:
		// This message is sent by the command to trigger a re-render after cursor change
		// No state change needed here
		
	case IncrementMsg:
		m.Count++
		m.Message = fmt.Sprintf("Count: %d. Press 'c' to toggle cursor, 'q' to quit", m.Count)
	}
	
	return m, nil
}

// View renders the model to a component tree
func (m CursorDemoModel) View() trmnl.Component {
	// Create title
	title := trmnl.NewText("Cursor Control Demo")
	title.Style.Border = true
	title.Style.RoundedBorder = true
	title.Style.BorderColor = trmnl.ColorBlue
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create status display
	status := fmt.Sprintf("Cursor Status: %s", func() string {
		if m.CursorVisible {
			return "VISIBLE"
		}
		return "HIDDEN"
	}())
	
	statusColor := trmnl.ColorGreen
	if m.CursorVisible {
		statusColor = trmnl.ColorYellow
	}
	
	statusDisplay := trmnl.NewText(status)
	statusDisplay.Style.Border = true
	statusDisplay.Style.BorderColor = statusColor
	statusDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create counter display
	counter := trmnl.NewText(fmt.Sprintf("Count: %d", m.Count))
	counter.Style.Border = true
	counter.Style.BorderColor = trmnl.ColorCyan
	counter.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create message display
	message := trmnl.NewText(m.Message)
	message.Style.Border = true
	message.Style.BorderColor = trmnl.ColorMagenta
	message.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create instructions
	instructions := trmnl.NewText("Instructions:\n  c : Toggle cursor visibility\n  + : Increment counter\n  q : Quit")
	instructions.Style.Border = true
	instructions.Style.BorderColor = trmnl.ColorBrightBlack
	instructions.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Layout everything
	return trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignCenter).
		SetPadding(trmnl.Box{Top: 2, Right: 2, Bottom: 2, Left: 2}).
		AddChild(title).
		AddChild(statusDisplay).
		AddChild(counter).
		AddChild(message).
		AddChild(instructions)
}

// Subscriptions returns no subscriptions for this demo
func (m CursorDemoModel) Subscriptions() []trmnl.Sub {
	return []trmnl.Sub{}
}

// toggleCursorCmd creates a command that toggles cursor visibility
func toggleCursorCmd(visible bool) trmnl.Cmd {
	return func() trmnl.Msg {
		// Note: In a real implementation, you might want to get access to the terminal
		// For now, this just sends a message to trigger a re-render
		// The actual cursor control happens via the Program's built-in management
		return ToggleCursorMsg{}
	}
}

func main() {
	// Create initial model
	initialModel := CursorDemoModel{
		Count:         0,
		CursorVisible: false,
	}
	
	fmt.Println("=== Cursor Control Demo ===")
	fmt.Println("This demo starts with cursor HIDDEN (default for TUI apps)")
	fmt.Println("Press any key to start...")
	
	// Wait for user input to start
	var input string
	fmt.Scanln(&input)
	
	// Create and run program with default cursor hidden
	program := trmnl.NewProgram(initialModel)
	if err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
	
	fmt.Println("\n=== Now trying with cursor VISIBLE ===")
	fmt.Println("Press any key to start cursor visible demo...")
	fmt.Scanln(&input)
	
	// Create and run program with cursor visible
	initialModel.CursorVisible = true
	initialModel.Message = "Cursor is VISIBLE by default. Press 'c' to toggle, 'q' to quit"
	program2 := trmnl.NewProgram(initialModel).WithCursorHidden(false)
	if err := program2.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
	
	fmt.Println("Demo complete!")
}