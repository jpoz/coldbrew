package main

import (
	"fmt"
	"time"

	"github.com/jpoz/trmnl"
)

// DemoModel shows difference between old (flickering) and new (smooth) rendering
type DemoModel struct {
	counter     int
	useOldStyle bool
}

func (m DemoModel) Init() (trmnl.Model, trmnl.Cmd) {
	return m, trmnl.Tick(100*time.Millisecond, func(t time.Time) trmnl.Msg {
		return TickMsg{}
	})
}

type TickMsg struct{}

func (m DemoModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.counter++
		return m, trmnl.Tick(100*time.Millisecond, func(t time.Time) trmnl.Msg {
			return TickMsg{}
		})
	case trmnl.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, trmnl.Quit()
		case "s":
			// Toggle style for comparison (if we had old implementation to compare)
			m.useOldStyle = !m.useOldStyle
		}
	}
	return m, nil
}

func (m DemoModel) View() trmnl.Component {
	// Dynamic content that changes every update
	title := trmnl.NewText("🚀 Diff-Based Rendering Demo")
	title.Style.Border = true
	title.Style.BorderColor = trmnl.ColorCyan
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	counterDisplay := trmnl.NewText(fmt.Sprintf("Frame: %d", m.counter))
	counterDisplay.Style.Border = true
	counterDisplay.Style.BorderColor = trmnl.ColorGreen
	counterDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 3, Bottom: 1, Left: 3}
	
	// Rapidly changing time
	timeDisplay := trmnl.NewText(fmt.Sprintf("Time: %s", time.Now().Format("15:04:05.000")))
	timeDisplay.Style.Border = true
	timeDisplay.Style.BorderColor = trmnl.ColorYellow
	timeDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Some content that doesn't change (to test diff efficiency)
	staticContent := trmnl.NewText("This text stays the same\nIt should not flicker\nEven during rapid updates")
	staticContent.Style.Border = true
	staticContent.Style.BorderColor = trmnl.ColorBlue
	staticContent.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Instructions
	instructions := trmnl.NewText("Notice: No screen clearing or flicker!\nOnly changed lines are updated.\nPress 'q' to quit.")
	instructions.Style.Border = true
	instructions.Style.BorderColor = trmnl.ColorMagenta
	instructions.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	container := trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignCenter).
		SetPadding(trmnl.Box{Top: 2, Right: 4, Bottom: 2, Left: 4}).
		AddChild(title).
		AddChild(counterDisplay).
		AddChild(timeDisplay).
		AddChild(staticContent).
		AddChild(instructions)
		
	return container
}

func (m DemoModel) Subscriptions() []trmnl.Sub {
	return nil
}

func main() {
	fmt.Println("🎯 Diff-Based Rendering Test")
	fmt.Println("This demo shows the new flicker-free rendering system.")
	fmt.Println("Before: Full screen clear + redraw (caused flicker)")
	fmt.Println("After: Only changed lines are updated (smooth)")
	fmt.Println()
	fmt.Print("Press Enter to start...")
	fmt.Scanln()
	
	model := DemoModel{counter: 0}
	
	program := trmnl.NewProgram(model).
		WithCursorHidden(true).
		WithRawMode(true)
	
	if err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println("Demo completed. The rendering should have been smooth with no flicker!")
}