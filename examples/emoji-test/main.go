package main

import (
	"fmt"

	"github.com/jpoz/trmnl"
)

// EmojiTestModel tests emoji rendering
type EmojiTestModel struct{}

func (m EmojiTestModel) Init() (trmnl.Model, trmnl.Cmd) {
	return m, nil
}

func (m EmojiTestModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case trmnl.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, trmnl.Quit()
		}
	}
	return m, nil
}

func (m EmojiTestModel) View() trmnl.Component {
	// Test various emojis
	title := trmnl.NewText("🧪 Emoji Test")
	title.Style.Border = true
	title.Style.BorderColor = trmnl.ColorCyan
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	simpleEmojis := trmnl.NewText("😀 😂 🤔 😍 🥺 🎉")
	simpleEmojis.Style.Border = true
	simpleEmojis.Style.BorderColor = trmnl.ColorYellow
	simpleEmojis.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	complexEmojis := trmnl.NewText("🔥 ⚡ 🌈 🦄 🎯 🚀")
	complexEmojis.Style.Border = true
	complexEmojis.Style.BorderColor = trmnl.ColorGreen
	complexEmojis.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	mixedContent := trmnl.NewText("Text before 🎭 emoji middle 🎪 text after")
	mixedContent.Style.Border = true
	mixedContent.Style.BorderColor = trmnl.ColorMagenta
	mixedContent.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	multilineEmojis := trmnl.NewText("Line 1: 📝 Writing\nLine 2: 🎨 Art\nLine 3: 🎵 Music")
	multilineEmojis.Style.Border = true
	multilineEmojis.Style.BorderColor = trmnl.ColorBlue
	multilineEmojis.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	instructions := trmnl.NewText("Press 'q' to quit\nCheck for emoji alignment issues")
	instructions.Style.Border = true
	instructions.Style.BorderColor = trmnl.ColorWhite
	instructions.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	container := trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignCenter).
		SetPadding(trmnl.Box{Top: 2, Right: 4, Bottom: 2, Left: 4}).
		AddChild(title).
		AddChild(simpleEmojis).
		AddChild(complexEmojis).
		AddChild(mixedContent).
		AddChild(multilineEmojis).
		AddChild(instructions)
		
	return container
}

func (m EmojiTestModel) Subscriptions() []trmnl.Sub {
	return nil
}

func main() {
	model := EmojiTestModel{}
	
	program := trmnl.NewProgram(model).
		WithCursorHidden(true).
		WithRawMode(true)
	
	if err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}