package trmnl

import (
	"fmt"
	"strings"
)

// BuildText creates a simple text string with optional styling
func BuildText(content string) string {
	return content
}

// BuildBorder creates a bordered text block
func BuildBorder(content, title string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return ""
	}
	
	// Calculate the width needed
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	
	// Add title if provided
	if title != "" && len(title) > maxWidth {
		maxWidth = len(title)
	}
	
	// Add padding for border
	width := maxWidth + 2
	
	var result strings.Builder
	
	// Top border
	if title != "" {
		titlePadding := (width - len(title) - 2) / 2
		result.WriteString("┌")
		result.WriteString(strings.Repeat("─", titlePadding))
		result.WriteString(" " + title + " ")
		result.WriteString(strings.Repeat("─", width-titlePadding-len(title)-3))
		result.WriteString("┐\n")
	} else {
		result.WriteString("┌")
		result.WriteString(strings.Repeat("─", width-2))
		result.WriteString("┐\n")
	}
	
	// Content lines
	for _, line := range lines {
		result.WriteString("│ ")
		result.WriteString(line)
		result.WriteString(strings.Repeat(" ", maxWidth-len(line)))
		result.WriteString(" │\n")
	}
	
	// Bottom border
	result.WriteString("└")
	result.WriteString(strings.Repeat("─", width-2))
	result.WriteString("┘")
	
	return result.String()
}

// BuildRoundedBorder creates a bordered text block with rounded corners
func BuildRoundedBorder(content, title string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return ""
	}
	
	// Calculate the width needed
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	
	// Add title if provided
	if title != "" && len(title) > maxWidth {
		maxWidth = len(title)
	}
	
	// Add padding for border
	width := maxWidth + 2
	
	var result strings.Builder
	
	// Top border
	if title != "" {
		titlePadding := (width - len(title) - 2) / 2
		result.WriteString("╭")
		result.WriteString(strings.Repeat("─", titlePadding))
		result.WriteString(" " + title + " ")
		result.WriteString(strings.Repeat("─", width-titlePadding-len(title)-3))
		result.WriteString("╮\n")
	} else {
		result.WriteString("╭")
		result.WriteString(strings.Repeat("─", width-2))
		result.WriteString("╮\n")
	}
	
	// Content lines
	for _, line := range lines {
		result.WriteString("│ ")
		result.WriteString(line)
		result.WriteString(strings.Repeat(" ", maxWidth-len(line)))
		result.WriteString(" │\n")
	}
	
	// Bottom border
	result.WriteString("╰")
	result.WriteString(strings.Repeat("─", width-2))
	result.WriteString("╯")
	
	return result.String()
}

// BuildFlex creates a flexible layout string combining multiple content blocks
func BuildFlex(direction Direction, items ...string) string {
	if len(items) == 0 {
		return ""
	}
	
	if direction == Row {
		return buildFlexRow(items...)
	} else {
		return buildFlexColumn(items...)
	}
}

// buildFlexRow combines items horizontally
func buildFlexRow(items ...string) string {
	if len(items) == 0 {
		return ""
	}
	
	// Split each item into lines
	itemLines := make([][]string, len(items))
	maxHeight := 0
	
	for i, item := range items {
		lines := strings.Split(item, "\n")
		itemLines[i] = lines
		if len(lines) > maxHeight {
			maxHeight = len(lines)
		}
	}
	
	// Pad all items to same height
	for i := range itemLines {
		for len(itemLines[i]) < maxHeight {
			itemLines[i] = append(itemLines[i], "")
		}
	}
	
	var result strings.Builder
	for row := 0; row < maxHeight; row++ {
		for i, lines := range itemLines {
			if i > 0 {
				result.WriteString(" ") // Space between items
			}
			result.WriteString(lines[row])
		}
		if row < maxHeight-1 {
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

// buildFlexColumn combines items vertically
func buildFlexColumn(items ...string) string {
	var result strings.Builder
	for i, item := range items {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(item)
	}
	return result.String()
}

// BuildCenter centers content within a specified width
func BuildCenter(content string, width int) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder
	
	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		
		lineWidth := len(line)
		if lineWidth >= width {
			result.WriteString(line)
		} else {
			padding := (width - lineWidth) / 2
			result.WriteString(strings.Repeat(" ", padding))
			result.WriteString(line)
			result.WriteString(strings.Repeat(" ", width-lineWidth-padding))
		}
	}
	
	return result.String()
}

// BuildPadding adds padding around content
func BuildPadding(content string, padding Box) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder
	
	// Top padding
	for i := 0; i < padding.Top; i++ {
		result.WriteString("\n")
	}
	
	// Content with left/right padding
	for i, line := range lines {
		if i > 0 || padding.Top > 0 {
			result.WriteString("\n")
		}
		result.WriteString(strings.Repeat(" ", padding.Left))
		result.WriteString(line)
		result.WriteString(strings.Repeat(" ", padding.Right))
	}
	
	// Bottom padding
	for i := 0; i < padding.Bottom; i++ {
		result.WriteString("\n")
	}
	
	return result.String()
}

// BuildProgress creates a progress bar string
func BuildProgress(current, total int, width int, title string) string {
	if width < 4 {
		width = 4
	}
	
	percentage := float64(current) / float64(total)
	if percentage > 1.0 {
		percentage = 1.0
	}
	
	filled := int(float64(width-2) * percentage)
	empty := width - 2 - filled
	
	var result strings.Builder
	
	if title != "" {
		result.WriteString(title + "\n")
	}
	
	result.WriteString("[")
	result.WriteString(strings.Repeat("█", filled))
	result.WriteString(strings.Repeat("░", empty))
	result.WriteString("]")
	result.WriteString(fmt.Sprintf(" %d/%d (%.0f%%)", current, total, percentage*100))
	
	return result.String()
}