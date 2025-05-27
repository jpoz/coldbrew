package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	brew "github.com/jpoz/coldbrew"
)

type model struct {
	input   textinput.Model
	history []string
	width   int
	height  int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type something..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40

	histroy := []string{}

	for i := range 50 {
		histroy = append(histroy, fmt.Sprintf("History entry %d", i+1))
	}

	return model{
		input:   ti,
		history: histroy,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", "ctrl+j":
			m.history = append([]string{m.input.Value()}, m.history...)
			m.input.SetValue("")
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	// History section
	if len(m.history) == 0 {
		b.WriteString("No entries yet\n\n")
	} else {
		for i := len(m.history) - 1; i >= 0; i-- {
			b.WriteString(m.history[i] + "\n")
		}
		b.WriteString("\n")
	}

	// Input box with border
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Width(m.input.Width + 2)

	b.WriteString(boxStyle.Render(m.input.View()))

	return b.String()
}

func main() {
	p := brew.NewProgram(initialModel())
	if err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
