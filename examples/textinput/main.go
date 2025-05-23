package main

import (
	"fmt"
	"strings"

	"github.com/jpoz/trmnl"
)

func main() {
	terminal := trmnl.NewTerminal()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Example 2: Centered Layout")
	fmt.Println(strings.Repeat("=", 80))

	// Example 2: Centered layout
	centerBox := trmnl.NewText("Centered\nContent")
	centerBox.Style.Border = true
	centerBox.Style.Padding = trmnl.Box{Top: 2, Right: 4, Bottom: 2, Left: 4}

	centeredContainer := trmnl.NewFlexContainer().
		SetDirection(trmnl.Row).
		SetJustify(trmnl.JustifyCenter).
		SetAlign(trmnl.AlignCenter).
		SetBorder(true).
		AddChild(centerBox)

	terminal.Render(centeredContainer, trmnl.Size{Width: 60, Height: 15})

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Example 3: Space Between Layout")
	fmt.Println(strings.Repeat("=", 80))

	// Example 3: Space between layout
	box1 := trmnl.NewText("Box 1")
	box1.Style.Border = true
	box1.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}

	box2 := trmnl.NewText("Box 2")
	box2.Style.Border = true
	box2.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}

	box3 := trmnl.NewText("Box 3")
	box3.Style.Border = true
	box3.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}

	spaceBetweenContainer := trmnl.NewFlexContainer().
		SetDirection(trmnl.Row).
		SetJustify(trmnl.JustifySpaceBetween).
		SetAlign(trmnl.AlignCenter).
		SetBorder(true).
		SetPadding(trmnl.Box{Top: 2, Right: 2, Bottom: 2, Left: 2}).
		AddChild(box1).
		AddChild(box2).
		AddChild(box3)

	terminal.Render(spaceBetweenContainer, trmnl.Size{Width: 70, Height: 10})
}
