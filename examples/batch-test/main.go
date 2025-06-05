package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	brew "github.com/jpoz/coldbrew"
)

type tickMsg time.Time

type model struct {
	count int
}

func (m model) Init() tea.Cmd {
	// Test tea.Batch - should receive 3 tick messages
	return tea.Batch(
		tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) }),
		tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) }),
		tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) }),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			return m, tea.Quit
		}
		// Test another batch
		return m, tea.Batch(
			tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) }),
			tea.Tick(75*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) }),
		)

	case tickMsg:
		m.count++
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Batch test - received %d tick messages.\n\nPress any key to send more batched commands.\nPress 'q' to quit.\n", m.count)
}

func main() {
	p := brew.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}