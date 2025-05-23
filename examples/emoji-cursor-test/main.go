package main

import (
	"fmt"
	"time"

	"github.com/jpoz/trmnl"
)

// EmojiCursorTestModel tests cursor positioning with emojis in diff-based rendering
type EmojiCursorTestModel struct {
	counter int
	phase   int
}

func (m EmojiCursorTestModel) Init() (trmnl.Model, trmnl.Cmd) {
	return m, trmnl.Tick(300*time.Millisecond, func(t time.Time) trmnl.Msg {
		return TickMsg{}
	})
}

type TickMsg struct{}

func (m EmojiCursorTestModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.counter++
		m.phase = (m.phase + 1) % 4
		return m, trmnl.Tick(300*time.Millisecond, func(t time.Time) trmnl.Msg {
			return TickMsg{}
		})
	case trmnl.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, trmnl.Quit()
		}
	}
	return m, nil
}

func (m EmojiCursorTestModel) View() trmnl.Component {
	title := trmnl.NewText("🎯 Emoji Cursor Position Test")
	title.Style.Border = true
	title.Style.BorderColor = trmnl.ColorCyan
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Test changing content with emojis - only some lines change
	line1 := trmnl.NewText("Static: 🎨 🎵 🎮 (this line never changes)")
	line1.Style.Border = true
	line1.Style.BorderColor = trmnl.ColorGreen
	line1.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// This line changes - should update cursor position correctly
	phaseEmojis := []string{"🔴", "🟡", "🟢", "🔵"}
	line2 := trmnl.NewText(fmt.Sprintf("Dynamic: %s %d (changes every frame)", phaseEmojis[m.phase], m.counter))
	line2.Style.Border = true
	line2.Style.BorderColor = trmnl.ColorYellow
	line2.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Another static line
	line3 := trmnl.NewText("Static: 🦄 🌈 ⭐ (this line never changes)")
	line3.Style.Border = true
	line3.Style.BorderColor = trmnl.ColorBlue
	line3.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// This line also changes - test multiple dynamic lines
	spinnerEmojis := []string{"⏰", "⏱️", "⏲️", "⏰"}
	line4 := trmnl.NewText(fmt.Sprintf("Clock: %s %s", spinnerEmojis[m.counter%4], time.Now().Format("15:04:05")))
	line4.Style.Border = true
	line4.Style.BorderColor = trmnl.ColorMagenta
	line4.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Final static line
	line5 := trmnl.NewText("Static: 🔥 ⚡ 🚀 (this line never changes)")
	line5.Style.Border = true
	line5.Style.BorderColor = trmnl.ColorWhite
	line5.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	instructions := trmnl.NewText("Watch for emoji alignment issues in diff updates\nPress 'q' to quit")
	instructions.Style.Border = true
	instructions.Style.BorderColor = trmnl.ColorRed
	instructions.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	container := trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignCenter).
		SetPadding(trmnl.Box{Top: 2, Right: 4, Bottom: 2, Left: 4}).
		AddChild(title).
		AddChild(line1).
		AddChild(line2).
		AddChild(line3).
		AddChild(line4).
		AddChild(line5).
		AddChild(instructions)
		
	return container
}

func (m EmojiCursorTestModel) Subscriptions() []trmnl.Sub {
	return nil
}

func main() {
	fmt.Println("🧪 Emoji Cursor Position Test")
	fmt.Println("This tests diff-based rendering with emojis.")
	fmt.Println("Expected behavior:")
	fmt.Println("  - Static lines should never flicker or shift")
	fmt.Println("  - Only dynamic lines should update")
	fmt.Println("  - Emoji alignment should be correct")
	fmt.Println("  - Cursor positioning should be accurate")
	fmt.Println()
	fmt.Print("Press Enter to start...")
	fmt.Scanln()
	
	model := EmojiCursorTestModel{counter: 0, phase: 0}
	
	program := trmnl.NewProgram(model).
		WithCursorHidden(true).
		WithRawMode(true)
	
	if err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println("Test completed. Emojis should have rendered correctly with proper diff updates!")
}