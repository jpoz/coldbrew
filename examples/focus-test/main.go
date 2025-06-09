package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	brew "github.com/jpoz/coldbrew"
)

type model struct {
	focused     bool
	lastEvent   string
	eventCount  int
	initialized bool
}

func (m model) Init() tea.Cmd {
	// Enable focus reporting when the program starts (using bubbletea's command for full compatibility)
	return tea.EnableReportFocus
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			// Disable focus reporting before quitting (using bubbletea's command)
			return m, tea.Batch(
				tea.DisableReportFocus,
				tea.Quit,
			)
		}
		if s := msg.String(); s == "r" {
			// Reset counters
			m.eventCount = 0
			m.lastEvent = "reset"
			return m, nil
		}

	case tea.FocusMsg:
		m.focused = true
		m.lastEvent = fmt.Sprintf("Focus gained at %s", time.Now().Format("15:04:05"))
		m.eventCount++
		m.initialized = true
		return m, nil

	case tea.BlurMsg:
		m.focused = false
		m.lastEvent = fmt.Sprintf("Focus lost at %s", time.Now().Format("15:04:05"))
		m.eventCount++
		m.initialized = true
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	status := "Unknown"
	if m.initialized {
		if m.focused {
			status = "Focused 🟢"
		} else {
			status = "Blurred 🔴"
		}
	} else {
		status = "Waiting for first focus event..."
	}

	return fmt.Sprintf(`Focus/Blur Test

Terminal Status: %s
Last Event: %s
Event Count: %d

Instructions:
- Click in/out of this terminal window to test focus/blur
- The status should change when the terminal gains/loses focus
- Press 'r' to reset counters
- Press 'q' to quit

Note: Focus reporting must be supported by your terminal.
`, status, m.lastEvent, m.eventCount)
}

func main() {
	p := brew.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}