package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jpoz/trmnl"
)

type model struct {
	message string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		default:
			m.message = fmt.Sprintf("Pressed: %s (press 'q' to quit)", msg.String())
		}
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Simple Quit Test\n\n%s\n\nPress 'q' to quit with tea.Quit", m.message)
}

func main() {
	m := model{message: "Press any key..."}
	p := trmnl.NewProgram(m).WithRawMode(true)
	
	if err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println("Successfully quit!")
}