package trmnl

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// FlexContainer implements flexbox-like layout
type FlexContainer struct {
	Children   []Component
	Direction  Direction
	Justify    Justify
	Align      Align
	Style      Style
	FlexGrow   map[Component]int
	FlexShrink map[Component]int
}

func NewFlexContainer() *FlexContainer {
	return &FlexContainer{
		Children:   make([]Component, 0),
		Direction:  Row,
		Justify:    JustifyStart,
		Align:      AlignStart,
		FlexGrow:   make(map[Component]int),
		FlexShrink: make(map[Component]int),
		Style: Style{
			BorderChar: '│',
			BgChar:     ' ',
		},
	}
}

func (fc *FlexContainer) AddChild(child Component) *FlexContainer {
	fc.Children = append(fc.Children, child)
	return fc
}

func (fc *FlexContainer) SetDirection(dir Direction) *FlexContainer {
	fc.Direction = dir
	return fc
}

func (fc *FlexContainer) SetJustify(justify Justify) *FlexContainer {
	fc.Justify = justify
	return fc
}

func (fc *FlexContainer) SetAlign(align Align) *FlexContainer {
	fc.Align = align
	return fc
}

func (fc *FlexContainer) SetBorder(border bool) *FlexContainer {
	fc.Style.Border = border
	return fc
}

func (fc *FlexContainer) SetRoundedBorder(rounded bool) *FlexContainer {
	fc.Style.RoundedBorder = rounded
	return fc
}

func (fc *FlexContainer) SetBorderColor(color Color) *FlexContainer {
	fc.Style.BorderColor = color
	return fc
}

func (fc *FlexContainer) SetPadding(padding Box) *FlexContainer {
	fc.Style.Padding = padding
	return fc
}

func (fc *FlexContainer) SetFlexGrow(child Component, grow int) *FlexContainer {
	fc.FlexGrow[child] = grow
	return fc
}

func (fc *FlexContainer) GetStyle() Style {
	return fc.Style
}

func (fc *FlexContainer) GetMinSize() Size {
	if len(fc.Children) == 0 {
		return Size{Width: 1, Height: 1}
	}

	var totalWidth, totalHeight int
	var maxWidth, maxHeight int

	for _, child := range fc.Children {
		childSize := child.GetMinSize()

		if fc.Direction == Row {
			totalWidth += childSize.Width
			if childSize.Height > maxHeight {
				maxHeight = childSize.Height
			}
		} else {
			totalHeight += childSize.Height
			if childSize.Width > maxWidth {
				maxWidth = childSize.Width
			}
		}
	}

	var width, height int
	if fc.Direction == Row {
		width = totalWidth
		height = maxHeight
	} else {
		width = maxWidth
		height = totalHeight
	}

	// Add container padding and border
	extraWidth := fc.Style.Padding.Left + fc.Style.Padding.Right
	extraHeight := fc.Style.Padding.Top + fc.Style.Padding.Bottom

	if fc.Style.Border {
		extraWidth += 2
		extraHeight += 2
	}

	return Size{
		Width:  width + extraWidth,
		Height: height + extraHeight,
	}
}

func (fc *FlexContainer) RenderToBuffer(size Size) *RenderBuffer {
	buffer := NewRenderBuffer(size.Height)

	// Initialize with background
	for i := range buffer.Lines {
		buffer.Lines[i] = strings.Repeat(string(fc.Style.BgChar), size.Width)
	}

	// Calculate content area
	contentWidth := size.Width
	contentHeight := size.Height
	startX := 0
	startY := 0

	if fc.Style.Border {
		contentWidth -= 2
		contentHeight -= 2
		startX = 1
		startY = 1

		// Draw border (plain text)
		borderH := string(fc.Style.BorderChar)
		borderV := "─"
		var corners []rune
		if fc.Style.RoundedBorder {
			corners = []rune("╭╮╰╯")
		} else {
			corners = []rune("┌┐└┘")
		}

		buffer.Lines[0] = string(corners[0]) + strings.Repeat(borderV, size.Width-2) + string(corners[1])
		buffer.Lines[size.Height-1] = string(corners[2]) + strings.Repeat(borderV, size.Width-2) + string(corners[3])

		for i := 1; i < size.Height-1; i++ {
			buffer.Lines[i] = borderH + buffer.Lines[i][1:len(buffer.Lines[i])-1] + borderH
		}

		// Add color information for borders
		if fc.Style.BorderColor != ColorDefault {
			// Color entire top and bottom borders
			buffer.AddColor(0, 0, size.Width, fc.Style.BorderColor)
			buffer.AddColor(size.Height-1, 0, size.Width, fc.Style.BorderColor)
			
			// Color side borders
			for i := 1; i < size.Height-1; i++ {
				buffer.AddColor(i, 0, 1, fc.Style.BorderColor)             // Left border
				buffer.AddColor(i, size.Width-1, size.Width, fc.Style.BorderColor) // Right border
			}
		}
	}

	contentWidth -= fc.Style.Padding.Left + fc.Style.Padding.Right
	contentHeight -= fc.Style.Padding.Top + fc.Style.Padding.Bottom
	contentStartX := startX + fc.Style.Padding.Left
	contentStartY := startY + fc.Style.Padding.Top

	if len(fc.Children) == 0 {
		return buffer
	}

	// Calculate child sizes and positions (same layout logic as before)
	childSizes := make([]Size, len(fc.Children))
	childPositions := make([]Position, len(fc.Children))

	// Get minimum sizes
	totalMinSize := 0
	for i, child := range fc.Children {
		minSize := child.GetMinSize()
		childSizes[i] = minSize

		if fc.Direction == Row {
			totalMinSize += minSize.Width
		} else {
			totalMinSize += minSize.Height
		}
	}

	// Calculate available space and distribute
	var availableSpace int
	if fc.Direction == Row {
		availableSpace = contentWidth
	} else {
		availableSpace = contentHeight
	}

	extraSpace := availableSpace - totalMinSize
	if extraSpace > 0 {
		// Distribute extra space based on flex-grow
		totalGrow := 0
		for _, child := range fc.Children {
			if grow, exists := fc.FlexGrow[child]; exists {
				totalGrow += grow
			}
		}

		if totalGrow > 0 {
			for i, child := range fc.Children {
				if grow, exists := fc.FlexGrow[child]; exists && grow > 0 {
					extraForChild := (extraSpace * grow) / totalGrow
					if fc.Direction == Row {
						childSizes[i].Width += extraForChild
					} else {
						childSizes[i].Height += extraForChild
					}
				}
			}
		}
	}

	// Position children based on justify-content
	currentPos := 0
	spacing := 0

	if fc.Justify == JustifyCenter {
		remainingSpace := availableSpace
		for _, size := range childSizes {
			if fc.Direction == Row {
				remainingSpace -= size.Width
			} else {
				remainingSpace -= size.Height
			}
		}
		currentPos = remainingSpace / 2
	} else if fc.Justify == JustifyEnd {
		remainingSpace := availableSpace
		for _, size := range childSizes {
			if fc.Direction == Row {
				remainingSpace -= size.Width
			} else {
				remainingSpace -= size.Height
			}
		}
		currentPos = remainingSpace
	} else if fc.Justify == JustifySpaceBetween && len(fc.Children) > 1 {
		remainingSpace := availableSpace
		for _, size := range childSizes {
			if fc.Direction == Row {
				remainingSpace -= size.Width
			} else {
				remainingSpace -= size.Height
			}
		}
		spacing = remainingSpace / (len(fc.Children) - 1)
	} else if fc.Justify == JustifySpaceAround {
		remainingSpace := availableSpace
		for _, size := range childSizes {
			if fc.Direction == Row {
				remainingSpace -= size.Width
			} else {
				remainingSpace -= size.Height
			}
		}
		spacing = remainingSpace / len(fc.Children)
		currentPos = spacing / 2
	}

	// Set positions
	for i := range childSizes {
		if fc.Direction == Row {
			childPositions[i].X = currentPos

			// Handle align-items for cross axis
			switch fc.Align {
			case AlignCenter:
				childPositions[i].Y = (contentHeight - childSizes[i].Height) / 2
			case AlignEnd:
				childPositions[i].Y = contentHeight - childSizes[i].Height
			case AlignStretch:
				childSizes[i].Height = contentHeight
				childPositions[i].Y = 0
			default: // AlignStart
				childPositions[i].Y = 0
			}

			currentPos += childSizes[i].Width + spacing
		} else {
			childPositions[i].Y = currentPos

			// Handle align-items for cross axis
			switch fc.Align {
			case AlignCenter:
				childPositions[i].X = (contentWidth - childSizes[i].Width) / 2
			case AlignEnd:
				childPositions[i].X = contentWidth - childSizes[i].Width
			case AlignStretch:
				childSizes[i].Width = contentWidth
				childPositions[i].X = 0
			default: // AlignStart
				childPositions[i].X = 0
			}

			currentPos += childSizes[i].Height + spacing
		}
	}

	// Render children using new buffer system
	for i, child := range fc.Children {
		if childSizes[i].Width <= 0 || childSizes[i].Height <= 0 {
			continue
		}

		childBuffer := child.RenderToBuffer(childSizes[i])
		childX := contentStartX + childPositions[i].X
		childY := contentStartY + childPositions[i].Y

		// Merge child buffer into parent buffer
		for j, line := range childBuffer.Lines {
			y := childY + j
			if y >= len(buffer.Lines) || y < 0 {
				continue
			}

			// Replace content in parent buffer at the correct visual position
			parentLine := buffer.Lines[y]
			parentRunes := []rune(parentLine)
			
			// Calculate where to insert the child line
			insertPos := 0
			currentVisualPos := 0
			
			// Find the character position that corresponds to childX visual position
			for insertPos < len(parentRunes) && currentVisualPos < childX {
				currentVisualPos += runewidth.RuneWidth(parentRunes[insertPos])
				if currentVisualPos <= childX {
					insertPos++
				}
			}
			
			// Convert child line to runes and calculate its visual width
			childRunes := []rune(line)
			childVisualWidth := runewidth.StringWidth(line)
			
			// Replace the section in parent buffer
			if insertPos >= 0 && insertPos < len(parentRunes) {
				// Calculate how many characters to replace based on visual width
				replaceEnd := insertPos
				replacedVisualWidth := 0
				
				for replaceEnd < len(parentRunes) && replacedVisualWidth < childVisualWidth {
					replacedVisualWidth += runewidth.RuneWidth(parentRunes[replaceEnd])
					replaceEnd++
				}
				
				// Build new line: before + child + after
				var newLine []rune
				if insertPos > 0 {
					newLine = append(newLine, parentRunes[:insertPos]...)
				}
				newLine = append(newLine, childRunes...)
				if replaceEnd < len(parentRunes) {
					newLine = append(newLine, parentRunes[replaceEnd:]...)
				}
				
				buffer.Lines[y] = string(newLine)
			}

			// Merge color information with visual position adjustment
			for _, colorInfo := range childBuffer.ColorMaps[j] {
				adjustedStart := childX + colorInfo.Start
				adjustedEnd := childX + colorInfo.End
				if adjustedStart < size.Width && adjustedEnd > 0 {
					// Clamp to valid range
					if adjustedStart < 0 {
						adjustedStart = 0
					}
					if adjustedEnd > size.Width {
						adjustedEnd = size.Width
					}
					buffer.AddColor(y, adjustedStart, adjustedEnd, colorInfo.Color)
				}
			}
		}
	}

	return buffer
}

func (fc *FlexContainer) Render(size Size) []string {
	buffer := fc.RenderToBuffer(size)
	return buffer.ApplyColors()
}

