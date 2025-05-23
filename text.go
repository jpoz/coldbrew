package trmnl

import (
	"strings"
	"unicode/utf8"
)

// Text is a simple text component
type Text struct {
	Content string
	Style   Style
}

func NewText(content string) *Text {
	return &Text{
		Content: content,
		Style: Style{
			BorderChar: '│',
			BgChar:     ' ',
		},
	}
}

func (t *Text) GetStyle() Style {
	return t.Style
}

func (t *Text) GetMinSize() Size {
	lines := strings.Split(t.Content, "\n")
	maxWidth := 0
	for _, line := range lines {
		width := utf8.RuneCountInString(line)
		if width > maxWidth {
			maxWidth = width
		}
	}

	// Add padding and border
	extraWidth := t.Style.Padding.Left + t.Style.Padding.Right
	extraHeight := t.Style.Padding.Top + t.Style.Padding.Bottom

	if t.Style.Border {
		extraWidth += 2
		extraHeight += 2
	}

	return Size{
		Width:  maxWidth + extraWidth,
		Height: len(lines) + extraHeight,
	}
}

func (t *Text) RenderToBuffer(size Size) *RenderBuffer {
	buffer := NewRenderBuffer(size.Height)
	
	lines := strings.Split(t.Content, "\n")

	// Calculate content area
	contentWidth := size.Width
	contentHeight := size.Height
	startX := 0
	startY := 0

	if t.Style.Border {
		contentWidth -= 2
		contentHeight -= 2
		startX = 1
		startY = 1
	}

	contentWidth -= t.Style.Padding.Left + t.Style.Padding.Right
	contentHeight -= t.Style.Padding.Top + t.Style.Padding.Bottom

	// Initialize with background
	for i := range buffer.Lines {
		buffer.Lines[i] = strings.Repeat(string(t.Style.BgChar), size.Width)
	}

	// Draw border
	if t.Style.Border {
		borderH := string(t.Style.BorderChar)
		borderV := "─"
		var corners []rune
		if t.Style.RoundedBorder {
			corners = []rune("╭╮╰╯")
		} else {
			corners = []rune("┌┐└┘")
		}

		// Top and bottom borders (plain text)
		buffer.Lines[0] = string(corners[0]) + strings.Repeat(borderV, size.Width-2) + string(corners[1])
		buffer.Lines[size.Height-1] = string(corners[2]) + strings.Repeat(borderV, size.Width-2) + string(corners[3])

		// Side borders (plain text)
		for i := 1; i < size.Height-1; i++ {
			buffer.Lines[i] = borderH + buffer.Lines[i][1:len(buffer.Lines[i])-1] + borderH
		}

		// Add color information for borders
		if t.Style.BorderColor != ColorDefault {
			// Color entire top and bottom borders
			buffer.AddColor(0, 0, size.Width, t.Style.BorderColor)
			buffer.AddColor(size.Height-1, 0, size.Width, t.Style.BorderColor)
			
			// Color side borders
			for i := 1; i < size.Height-1; i++ {
				buffer.AddColor(i, 0, 1, t.Style.BorderColor)             // Left border
				buffer.AddColor(i, size.Width-1, size.Width, t.Style.BorderColor) // Right border
			}
		}
	}

	// Draw content
	contentStartY := startY + t.Style.Padding.Top
	contentStartX := startX + t.Style.Padding.Left

	for i, line := range lines {
		if i >= contentHeight {
			break
		}

		y := contentStartY + i
		if y >= len(buffer.Lines) {
			break
		}

		// Truncate or pad line to fit content width
		displayLine := line
		if utf8.RuneCountInString(line) > contentWidth {
			displayLine = string([]rune(line)[:contentWidth])
		} else if utf8.RuneCountInString(line) < contentWidth {
			displayLine = line + strings.Repeat(string(t.Style.BgChar), contentWidth-utf8.RuneCountInString(line))
		}

		// Replace the content portion of the result line
		resultRunes := []rune(buffer.Lines[y])
		displayRunes := []rune(displayLine)

		for j, r := range displayRunes {
			if contentStartX+j < len(resultRunes) {
				resultRunes[contentStartX+j] = r
			}
		}
		buffer.Lines[y] = string(resultRunes)
	}

	return buffer
}

func (t *Text) Render(size Size) []string {
	buffer := t.RenderToBuffer(size)
	return buffer.ApplyColors()
}