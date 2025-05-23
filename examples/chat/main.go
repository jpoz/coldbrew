package main

import (
	"fmt"

	"github.com/jpoz/trmnl"
)

// ChatModel represents our chat application state
type ChatModel struct {
	Input string
	Count int
}

// Init initializes the model
func (m ChatModel) Init() (trmnl.Model, trmnl.Cmd) {
	return m, nil
}

// Update handles messages
func (m ChatModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case trmnl.KeyMsg:
		if msg.String() == "q" {
			return m, trmnl.Quit()
		}
		m.Input = m.Input + msg.String()
		m.Count++
	}
	return m, nil
}

// View renders the model
func (m ChatModel) View() trmnl.Component {
	text := fmt.Sprintf("'%s'", m.Input)
	component := trmnl.NewText(text)
	component.Style.Border = true
	component.Style.RoundedBorder = true
	component.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	return component
}

// Subscriptions returns no subscriptions
func (m ChatModel) Subscriptions() []trmnl.Sub {
	return []trmnl.Sub{}
}

func main() {
	model := ChatModel{}
	program := trmnl.NewProgram(model).WithRawMode(true)

	if err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
