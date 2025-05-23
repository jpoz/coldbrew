package main

import (
	"fmt"

	"github.com/jpoz/trmnl"
)

// SimpleModel for testing immediate input
type SimpleModel struct {
	LastKey string
	Count   int
}

// Init initializes the model
func (m SimpleModel) Init() (trmnl.Model, trmnl.Cmd) {
	return m, nil
}

// Update handles messages
func (m SimpleModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case trmnl.KeyMsg:
		if msg.String() == "q" {
			return m, trmnl.Quit()
		}
		m.LastKey = msg.String()
		m.Count++
	}
	return m, nil
}

// View renders the model
func (m SimpleModel) View() trmnl.Component {
	text := fmt.Sprintf("Last key: '%s' | Count: %d | Press 'q' to quit", m.LastKey, m.Count)
	component := trmnl.NewText(text)
	component.Style.Border = true
	component.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	return component
}

// Subscriptions returns no subscriptions
func (m SimpleModel) Subscriptions() []trmnl.Sub {
	return []trmnl.Sub{}
}

func main() {
	fmt.Println("Simple immediate input test. Raw mode enabled.")
	fmt.Println("Type any keys (no Enter needed), 'q' to quit.")
	fmt.Println()
	
	model := SimpleModel{}
	program := trmnl.NewProgram(model).WithRawMode(true)
	
	if err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}