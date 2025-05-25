package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	brew "github.com/jpoz/coldbrew"
)

// QuitTestModel tests different quit methods
type QuitTestModel struct {
	Message string
}

// Init initializes the model
func (m QuitTestModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m QuitTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.Message = "Sending tea.Quit command..."
			return m, tea.Quit // This should exit the program
		case "x":
			m.Message = "Sending custom QuitMsg..."
			return m, func() tea.Msg { return brew.QuitMsg{} } // Custom quit
		case "ctrl+c":
			m.Message = "Sending ctrl+c tea.Quit..."
			return m, tea.Quit
		default:
			m.Message = fmt.Sprintf("Key: '%s' | Press q/x/ctrl+c to quit", msg.String())
		}
	case tea.QuitMsg:
		m.Message = "Received tea.QuitMsg - should quit now!"
		return m, nil
	case brew.QuitMsg:
		m.Message = "Received custom QuitMsg - should quit now!"
		return m, nil
	}
	return m, nil
}

// View renders the model
func (m QuitTestModel) View() string {
	content := fmt.Sprintf("Quit Test Demo\n\n%s\n\nTest different quit methods:\n• q - tea.Quit command\n• x - custom QuitMsg\n• ctrl+c - standard quit", m.Message)
	
	// Simple border since BuildRoundedBorder might not be available
	border := "╭" + strings.Repeat("─", 50) + "╮\n"
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		border += "│ " + line + strings.Repeat(" ", 48-len(line)) + " │\n"
	}
	border += "╰" + strings.Repeat("─", 50) + "╯"
	
	return border
}

func main() {
	fmt.Println("=== Quit Command Test ===")
	fmt.Println("Testing different ways to quit the program:")
	fmt.Println("• q - should use tea.Quit command")
	fmt.Println("• x - should use custom QuitMsg")
	fmt.Println("• ctrl+c - should also use tea.Quit")
	fmt.Println()
	fmt.Println("Press any key to start...")
	
	var input string
	fmt.Scanln(&input)
	
	model := QuitTestModel{Message: "Press q, x, or ctrl+c to test quit methods"}
	program := brew.NewProgram(model).WithRawMode(true)
	
	if err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println("Program exited successfully!")
}