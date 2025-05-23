package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jpoz/trmnl"
)

// ImmediateInputModel represents our immediate input demo
type ImmediateInputModel struct {
	LastKey     string
	KeyCount    int
	KeyHistory  []string
	Message     string
	CurrentTime time.Time
	RawMode     bool
}

// Message types
type TickMsg struct {
	Time time.Time
}

// Init initializes the model
func (m ImmediateInputModel) Init() (trmnl.Model, trmnl.Cmd) {
	if m.RawMode {
		m.Message = "RAW MODE: Press any key (no Enter needed)! 'q' to quit, arrow keys supported"
	} else {
		m.Message = "LINE MODE: Type keys and press Enter. 'q' to quit"
	}
	m.CurrentTime = time.Now()
	m.KeyHistory = []string{}
	return m, nil
}

// Update handles messages and returns updated model
func (m ImmediateInputModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case trmnl.KeyMsg:
		// Handle special keys
		switch msg.String() {
		case "q":
			return m, trmnl.Quit()
		case "ctrl+c":
			return m, trmnl.Quit()
		case "ctrl+d":
			return m, trmnl.Quit()
		default:
			m.LastKey = msg.String()
			m.KeyCount++
			
			// Add to history (keep last 10)
			m.KeyHistory = append(m.KeyHistory, msg.String())
			if len(m.KeyHistory) > 10 {
				m.KeyHistory = m.KeyHistory[1:]
			}
			
			// Create descriptive message
			switch msg.String() {
			case "up", "down", "left", "right":
				m.Message = fmt.Sprintf("Arrow key pressed: %s", msg.String())
			case "space":
				m.Message = "Space key pressed!"
			case "enter":
				m.Message = "Enter key pressed!"
			case "tab":
				m.Message = "Tab key pressed!"
			case "backspace":
				m.Message = "Backspace key pressed!"
			case "escape":
				m.Message = "Escape key pressed!"
			default:
				if msg.Rune != 0 {
					m.Message = fmt.Sprintf("Character '%c' pressed (key: %s)", msg.Rune, msg.String())
				} else {
					m.Message = fmt.Sprintf("Special key pressed: %s", msg.String())
				}
			}
		}
		
	case TickMsg:
		m.CurrentTime = msg.Time
	}
	
	return m, nil
}

// View renders the model to a component tree
func (m ImmediateInputModel) View() trmnl.Component {
	// Create title
	title := trmnl.NewText("Immediate Input Demo")
	title.Style.Border = true
	title.Style.RoundedBorder = true
	title.Style.BorderColor = trmnl.ColorBlue
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create current time display
	timeDisplay := trmnl.NewText(fmt.Sprintf("Current Time: %s", m.CurrentTime.Format("15:04:05")))
	timeDisplay.Style.Border = true
	timeDisplay.Style.BorderColor = trmnl.ColorCyan
	timeDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create last key display
	lastKeyText := "None"
	if m.LastKey != "" {
		lastKeyText = m.LastKey
	}
	
	lastKeyDisplay := trmnl.NewText(fmt.Sprintf("Last Key: %s", lastKeyText))
	lastKeyDisplay.Style.Border = true
	lastKeyDisplay.Style.BorderColor = trmnl.ColorGreen
	lastKeyDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create key count display
	countDisplay := trmnl.NewText(fmt.Sprintf("Keys Pressed: %d", m.KeyCount))
	countDisplay.Style.Border = true
	countDisplay.Style.BorderColor = trmnl.ColorYellow
	countDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create key history display
	historyText := "Key History (last 10):"
	if len(m.KeyHistory) == 0 {
		historyText += "\n  (no keys pressed yet)"
	} else {
		historyText += "\n  " + strings.Join(m.KeyHistory, " → ")
	}
	
	historyDisplay := trmnl.NewText(historyText)
	historyDisplay.Style.Border = true
	historyDisplay.Style.BorderColor = trmnl.ColorMagenta
	historyDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create message display
	messageDisplay := trmnl.NewText(m.Message)
	messageDisplay.Style.Border = true
	messageDisplay.Style.BorderColor = trmnl.ColorBrightGreen
	messageDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create instructions
	instructions := trmnl.NewText("Instructions:\n• Press any key to see immediate response\n• Arrow keys, space, enter, etc. all work\n• 'q' or Ctrl+C to quit\n• No need to press Enter!")
	instructions.Style.Border = true
	instructions.Style.BorderColor = trmnl.ColorBrightBlack
	instructions.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Layout everything in two columns
	leftColumn := trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignStretch).
		AddChild(title).
		AddChild(timeDisplay).
		AddChild(lastKeyDisplay).
		AddChild(countDisplay)
	
	rightColumn := trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignStretch).
		AddChild(historyDisplay).
		AddChild(messageDisplay).
		AddChild(instructions)
	
	return trmnl.NewFlexContainer().
		SetDirection(trmnl.Row).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignStart).
		SetPadding(trmnl.Box{Top: 1, Right: 1, Bottom: 1, Left: 1}).
		AddChild(leftColumn).
		AddChild(rightColumn).
		SetFlexGrow(leftColumn, 1).
		SetFlexGrow(rightColumn, 1)
}

// Subscriptions returns the subscriptions for this model
func (m ImmediateInputModel) Subscriptions() []trmnl.Sub {
	return []trmnl.Sub{
		trmnl.Every(time.Second, func(t time.Time) trmnl.Msg {
			return TickMsg{Time: t}
		}),
	}
}

func main() {
	fmt.Println("=== Immediate Input Demo ===")
	fmt.Println("This demo shows immediate character input without pressing Enter.")
	fmt.Println("The terminal will be put into 'raw mode' for immediate response.")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("• Each keypress is processed immediately")
	fmt.Println("• Arrow keys are supported")
	fmt.Println("• Special keys (space, tab, etc.) are detected")
	fmt.Println("• Real-time clock updates")
	fmt.Println("• Key history tracking")
	fmt.Println()
	fmt.Println("Press any key to start...")
	
	// Wait for user input to start
	var input string
	fmt.Scanln(&input)
	
	// Create initial model
	initialModel := ImmediateInputModel{
		KeyCount:    0,
		KeyHistory:  []string{},
		CurrentTime: time.Now(),
		RawMode:     true, // We're requesting raw mode
	}
	
	// Create and run program with raw mode enabled
	program := trmnl.NewProgram(initialModel).
		WithRawMode(true).
		WithCursorHidden(true)
	
	if err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
	
	fmt.Println("\nDemo complete! Terminal has been restored to normal mode.")
}