package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	brew "github.com/jpoz/coldbrew"
)

type model struct {
	count int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			return m, tea.Quit
		}
		if s := msg.String(); s == " " {
			m.count++
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Counter: %d\n\nPress space to increment.\nPress 'q' to quit.\n", m.count)
}

func main() {
	p := brew.NewProgram(model{})
	
	// Run returns both the final model and any error, just like bubbletea
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	// Access the final state of the model after the program exits
	if m, ok := finalModel.(model); ok {
		fmt.Printf("Program exited with final count: %d\n", m.count)
	} else {
		fmt.Println("Program exited")
	}
}