package main

import (
	"fmt"
	"time"

	"github.com/jpoz/trmnl"
)

// EmojiDynamicModel tests dynamic emoji rendering with diff updates
type EmojiDynamicModel struct {
	counter int
	emojis  []string
}

func (m EmojiDynamicModel) Init() (trmnl.Model, trmnl.Cmd) {
	m.emojis = []string{"🔴", "🟠", "🟡", "🟢", "🔵", "🟣", "⚫", "⚪"}
	return m, trmnl.Tick(500*time.Millisecond, func(t time.Time) trmnl.Msg {
		return TickMsg{}
	})
}

type TickMsg struct{}

func (m EmojiDynamicModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.counter++
		return m, trmnl.Tick(500*time.Millisecond, func(t time.Time) trmnl.Msg {
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

func (m EmojiDynamicModel) View() trmnl.Component {
	title := trmnl.NewText("🔄 Dynamic Emoji Test")
	title.Style.Border = true
	title.Style.BorderColor = trmnl.ColorCyan
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Rotating emoji indicator
	currentEmoji := m.emojis[m.counter%len(m.emojis)]
	emojiCounter := trmnl.NewText(fmt.Sprintf("Current: %s (Frame %d)", currentEmoji, m.counter))
	emojiCounter.Style.Border = true
	emojiCounter.Style.BorderColor = trmnl.ColorYellow
	emojiCounter.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Progress bar with emoji
	progress := m.counter % 10
	progressBar := ""
	for i := 0; i < 10; i++ {
		if i <= progress {
			progressBar += "🟩"
		} else {
			progressBar += "⬜"
		}
	}
	
	progressDisplay := trmnl.NewText(fmt.Sprintf("Progress: %s", progressBar))
	progressDisplay.Style.Border = true
	progressDisplay.Style.BorderColor = trmnl.ColorGreen
	progressDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Clock with emoji
	now := time.Now()
	clockDisplay := trmnl.NewText(fmt.Sprintf("🕐 Time: %s", now.Format("15:04:05")))
	clockDisplay.Style.Border = true
	clockDisplay.Style.BorderColor = trmnl.ColorBlue
	clockDisplay.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Static emoji content (should not flicker)
	staticContent := trmnl.NewText("Static: 🎯 🎨 🎵 🎮\nThese should stay stable")
	staticContent.Style.Border = true
	staticContent.Style.BorderColor = trmnl.ColorMagenta
	staticContent.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	instructions := trmnl.NewText("🔍 Watch for emoji alignment issues\nPress 'q' to quit")
	instructions.Style.Border = true
	instructions.Style.BorderColor = trmnl.ColorWhite
	instructions.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	container := trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignCenter).
		SetPadding(trmnl.Box{Top: 2, Right: 4, Bottom: 2, Left: 4}).
		AddChild(title).
		AddChild(emojiCounter).
		AddChild(progressDisplay).
		AddChild(clockDisplay).
		AddChild(staticContent).
		AddChild(instructions)
		
	return container
}

func (m EmojiDynamicModel) Subscriptions() []trmnl.Sub {
	return nil
}

func main() {
	fmt.Println("🧪 Dynamic Emoji Rendering Test")
	fmt.Println("This tests emoji rendering with diff-based updates.")
	fmt.Println("Watch for:")
	fmt.Println("  - Emoji alignment issues")
	fmt.Println("  - Character width problems")
	fmt.Println("  - Cursor positioning errors")
	fmt.Println()
	fmt.Print("Press Enter to start...")
	fmt.Scanln()
	
	model := EmojiDynamicModel{counter: 0}
	
	program := trmnl.NewProgram(model).
		WithCursorHidden(true).
		WithRawMode(true)
	
	if err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println("Test completed. Check if emojis rendered correctly!")
}