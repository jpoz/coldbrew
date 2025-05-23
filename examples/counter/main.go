package main

import (
	"fmt"

	"github.com/jpoz/trmnl"
)

// CounterModel represents our application state
type CounterModel struct {
	Count   int
	Message string
}

// Message types
type IncrementMsg struct{}
type DecrementMsg struct{}
type ResetMsg struct{}

// Init initializes the model
func (m CounterModel) Init() (trmnl.Model, trmnl.Cmd) {
	m.Message = "Use +/- to change count, 'r' to reset, 'q' to quit"
	return m, nil
}

// Update handles messages and returns updated model
func (m CounterModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case trmnl.KeyMsg:
		switch msg.String() {
		case "q":
			return m, trmnl.Quit()
		case "+", "=":
			m.Count++
			m.Message = fmt.Sprintf("Count incremented to %d", m.Count)
		case "-":
			m.Count--
			m.Message = fmt.Sprintf("Count decremented to %d", m.Count)
		case "r":
			m.Count = 0
			m.Message = "Count reset to 0"
		case "b":
			// Test batch messages - increment 3 times
			m.Message = "Batch increment triggered!"
			return m, trmnl.Batch(IncrementMsg{}, IncrementMsg{}, IncrementMsg{})
		default:
			m.Message = fmt.Sprintf("Unknown key: %s. Use +/-/b/r, 'q' to quit", msg.String())
		}
		
	case IncrementMsg:
		m.Count++
		m.Message = fmt.Sprintf("Count incremented to %d", m.Count)
		
	case DecrementMsg:
		m.Count--
		m.Message = fmt.Sprintf("Count decremented to %d", m.Count)
		
	case ResetMsg:
		m.Count = 0
		m.Message = "Count reset to 0"
	}
	
	return m, nil
}

// View renders the model to a component tree
func (m CounterModel) View() trmnl.Component {
	// Create title
	title := trmnl.NewText("Counter Example")
	title.Style.Border = true
	title.Style.RoundedBorder = true
	title.Style.BorderColor = trmnl.ColorBlue
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create counter display
	counterText := fmt.Sprintf("Count: %d", m.Count)
	counter := trmnl.NewText(counterText)
	counter.Style.Border = true
	counter.Style.BorderColor = trmnl.ColorGreen
	counter.Style.Padding = trmnl.Box{Top: 2, Right: 4, Bottom: 2, Left: 4}
	
	// Create message display
	message := trmnl.NewText(m.Message)
	message.Style.Border = true
	message.Style.BorderColor = trmnl.ColorYellow
	message.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create instructions
	instructions := trmnl.NewText("Commands:\n  + : Increment\n  - : Decrement\n  b : Batch +3\n  r : Reset\n  q : Quit")
	instructions.Style.Border = true
	instructions.Style.BorderColor = trmnl.ColorCyan
	instructions.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Layout everything
	return trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignCenter).
		SetPadding(trmnl.Box{Top: 2, Right: 2, Bottom: 2, Left: 2}).
		AddChild(title).
		AddChild(counter).
		AddChild(message).
		AddChild(instructions)
}

// Subscriptions returns no subscriptions for the counter
func (m CounterModel) Subscriptions() []trmnl.Sub {
	return []trmnl.Sub{}
}

func main() {
	// Create initial model
	initialModel := CounterModel{
		Count: 0,
	}
	
	// Create and run program
	program := trmnl.NewProgram(initialModel)
	if err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}