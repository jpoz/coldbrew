package main

import (
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

	// Render to terminal
	terminal.Render(rootContainer, trmnl.Size{Width: 80, Height: 20})

	// fmt.Println("\n" + strings.Repeat("=", 80))
	// fmt.Println("Example 2: Centered Layout")
	// fmt.Println(strings.Repeat("=", 80))
	//
	// // Example 2: Centered layout
	// centerBox := trmnl.NewText("Centered\nContent")
	// centerBox.Style.Border = true
	// centerBox.Style.Padding = trmnl.Box{Top: 2, Right: 4, Bottom: 2, Left: 4}
	//
	// centeredContainer := trmnl.NewFlexContainer().
	// 	SetDirection(trmnl.Row).
	// 	SetJustify(trmnl.JustifyCenter).
	// 	SetAlign(trmnl.AlignCenter).
	// 	SetBorder(true).
	// 	AddChild(centerBox)
	//
	// terminal.Render(centeredContainer, trmnl.Size{Width: 60, Height: 15})
	//
	// fmt.Println("\n" + strings.Repeat("=", 80))
	// fmt.Println("Example 3: Space Between Layout")
	// fmt.Println(strings.Repeat("=", 80))
	//
	// // Example 3: Space between layout
	// box1 := trmnl.NewText("Box 1")
	// box1.Style.Border = true
	// box1.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	//
	// box2 := trmnl.NewText("Box 2")
	// box2.Style.Border = true
	// box2.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	//
	// box3 := trmnl.NewText("Box 3")
	// box3.Style.Border = true
	// box3.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	//
	// spaceBetweenContainer := trmnl.NewFlexContainer().
	// 	SetDirection(trmnl.Row).
	// 	SetJustify(trmnl.JustifySpaceBetween).
	// 	SetAlign(trmnl.AlignCenter).
	// 	SetBorder(true).
	// 	SetPadding(trmnl.Box{Top: 2, Right: 2, Bottom: 2, Left: 2}).
	// 	AddChild(box1).
	// 	AddChild(box2).
	// 	AddChild(box3)
	//
	// terminal.Render(spaceBetweenContainer, trmnl.Size{Width: 70, Height: 10})
}
