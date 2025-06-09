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
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			return m, tea.Quit
		}

	case tickMsg:
		m.count++
		// Continue ticking
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Counter: %d\n\nThis will automatically quit after 5 seconds.\nOr press 'q' to quit manually.\n", m.count)
}

func main() {
	p := brew.NewProgram(model{})
	
	// Start a goroutine that will quit the program after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Calling program.Quit() from outside...")
		p.Quit() // This demonstrates the Program.Quit() method
	}()
	
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	if m, ok := finalModel.(model); ok {
		fmt.Printf("Program exited after %d seconds\n", m.count)
	}
}