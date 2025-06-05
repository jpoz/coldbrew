package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	brew "github.com/jpoz/coldbrew"
)

type model struct {
	width  int
	height int
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
		// Trigger a manual window size check
		return m, brew.WindowSize()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 && m.height == 0 {
		return "Press any key to check window size. Resize the terminal to see automatic updates. Press 'q' to quit.\n"
	}
	
	return fmt.Sprintf("Terminal size: %dx%d\n\nPress any key to manually check size.\nResize the terminal to see automatic updates.\nPress 'q' to quit.\n", m.width, m.height)
}

func main() {
	p := brew.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}