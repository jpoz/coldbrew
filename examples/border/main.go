package main

import (
	"fmt"
	"strings"

	"github.com/jpoz/trmnl"
)

func main() {
	terminal := trmnl.NewTerminal()
	terminal.Clear()

	// Create some text components
	header := trmnl.NewText("Header Text")
	header.Style.Border = true
	header.Style.RoundedBorder = true
	header.Style.BorderColor = trmnl.ColorBlue
	header.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}

	sidebar := trmnl.NewText("Sidebar\nContent\nHere")
	sidebar.Style.Border = true
	sidebar.Style.RoundedBorder = true
	sidebar.Style.BorderColor = trmnl.ColorGreen
	sidebar.Style.Padding = trmnl.Box{Top: 1, Right: 1, Bottom: 1, Left: 1}

	main := trmnl.NewText("Main Content Area\nThis is the main content\nwith multiple lines\nof text content.")
	main.Style.Border = true
	main.Style.RoundedBorder = true
	main.Style.BorderColor = trmnl.ColorCyan
	main.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}

	footer := trmnl.NewText("Footer - Status Bar")
	footer.Style.Border = true
	footer.Style.RoundedBorder = true
	footer.Style.BorderColor = trmnl.ColorYellow
	footer.Style.Padding = trmnl.Box{Top: 0, Right: 2, Bottom: 0, Left: 2}

	// Create layout structure
	bodyContainer := trmnl.NewFlexContainer().
		SetDirection(trmnl.Row).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignStretch).
		AddChild(sidebar).
		AddChild(main).
		SetFlexGrow(main, 1) // Main content takes extra space

	rootContainer := trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignStretch).
		SetBorder(true).
		SetRoundedBorder(true).
		SetBorderColor(trmnl.ColorMagenta).
		SetPadding(trmnl.Box{Top: 1, Right: 1, Bottom: 1, Left: 1}).
		AddChild(header).
		AddChild(bodyContainer).
		AddChild(footer).
		SetFlexGrow(bodyContainer, 1) // Body takes extra vertical space

	// Render to terminal using responsive methods
	fmt.Println("Example 1: Full Terminal Width Layout")
	fmt.Println(strings.Repeat("=", 80))
	terminal.RenderFullWidth(rootContainer, 30)
}
