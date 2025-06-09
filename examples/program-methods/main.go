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
	count     int
	startTime time.Time
}

func (m model) Init() tea.Cmd {
	m.startTime = time.Now()
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
	elapsed := time.Since(m.startTime).Seconds()
	return fmt.Sprintf(`Program Methods Demo

Counter: %d
Elapsed: %.1f seconds

This demonstrates all Program methods:
- Send() - sending messages
- Quit() - graceful quit from outside  
- Kill() - immediate termination
- Wait() - wait for program to finish
- Run() - returns (Model, error)

Press 'q' to quit, or wait for external quit...
`, m.count, elapsed)
}

func main() {
	fmt.Println("Starting program methods demo...")
	
	p := brew.NewProgram(model{startTime: time.Now()})
	
	// Start a goroutine that demonstrates external control
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("\n[External] Sending a custom message...")
		p.Send(tickMsg(time.Now()))
		
		time.Sleep(2 * time.Second)
		fmt.Println("\n[External] Calling program.Quit()...")
		p.Quit()
	}()
	
	// Start another goroutine that waits for the program to finish
	go func() {
		fmt.Println("[Waiter] Starting to wait for program completion...")
		p.Wait()
		fmt.Println("[Waiter] Program has finished!")
	}()
	
	// Run the program and get the final model
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	if m, ok := finalModel.(model); ok {
		elapsed := time.Since(m.startTime).Seconds()
		fmt.Printf("\n=== Program Results ===\n")
		fmt.Printf("Final count: %d\n", m.count)
		fmt.Printf("Total runtime: %.1f seconds\n", elapsed)
		fmt.Printf("Program completed successfully!\n")
	}
}