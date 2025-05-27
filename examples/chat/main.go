package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	brew "github.com/jpoz/coldbrew"
)

func main() {
	p := brew.NewProgram(model(5))
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model int

func (m model) Init() tea.Cmd {
	return tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		}

	case tickMsg:
		m--
		if m <= 0 {
			return m, tea.Quit
		}
		return m, tick
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("line 1\nline 2\nline 3\n")
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}
