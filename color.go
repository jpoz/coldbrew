package trmnl

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// ColorInfo stores color information for a specific position range in a string
type ColorInfo struct {
	Start int   // Start position (visual position, not byte position)
	End   int   // End position (visual position, not byte position)
	Color Color // Color to apply
}

// ColorMap stores color information for a line
type ColorMap []ColorInfo

// RenderBuffer represents a buffer with separate content and color information
type RenderBuffer struct {
	Lines      []string     // Plain text content for layout calculations
	ColorMaps  []ColorMap   // Color information for each line
}

// NewRenderBuffer creates a new render buffer
func NewRenderBuffer(height int) *RenderBuffer {
	return &RenderBuffer{
		Lines:     make([]string, height),
		ColorMaps: make([]ColorMap, height),
	}
}

// SetLine sets the content of a line
func (rb *RenderBuffer) SetLine(lineIndex int, content string) {
	if lineIndex >= 0 && lineIndex < len(rb.Lines) {
		rb.Lines[lineIndex] = content
	}
}

// AddColor adds color information to a specific range in a line
func (rb *RenderBuffer) AddColor(lineIndex, start, end int, color Color) {
	if lineIndex >= 0 && lineIndex < len(rb.ColorMaps) && color != ColorDefault {
		rb.ColorMaps[lineIndex] = append(rb.ColorMaps[lineIndex], ColorInfo{
			Start: start,
			End:   end,
			Color: color,
		})
	}
}

// ApplyColors applies all color information to produce final colored strings
func (rb *RenderBuffer) ApplyColors() []string {
	result := make([]string, len(rb.Lines))
	
	for i, line := range rb.Lines {
		result[i] = rb.applyColorsToLine(line, rb.ColorMaps[i])
	}
	
	return result
}

// applyColorsToLine applies color map to a single line
func (rb *RenderBuffer) applyColorsToLine(line string, colorMap ColorMap) string {
	if len(colorMap) == 0 {
		return line
	}
	
	runes := []rune(line)
	var result strings.Builder
	
	i := 0
	for i < len(runes) {
		// Find if current position has color
		var currentColor Color = ColorDefault
		for _, colorInfo := range colorMap {
			if i >= colorInfo.Start && i < colorInfo.End {
				currentColor = colorInfo.Color
				break
			}
		}
		
		if currentColor != ColorDefault {
			// Find the end of this color section
			end := i + 1
			for end < len(runes) {
				stillInColor := false
				for _, colorInfo := range colorMap {
					if end >= colorInfo.Start && end < colorInfo.End && colorInfo.Color == currentColor {
						stillInColor = true
						break
					}
				}
				if !stillInColor {
					break
				}
				end++
			}
			
			// Apply color to this section
			result.WriteString(string(currentColor))
			result.WriteString(string(runes[i:end]))
			result.WriteString(string(ColorReset))
			
			i = end
		} else {
			result.WriteRune(runes[i])
			i++
		}
	}
	
	return result.String()
}

// VisualWidth calculates the visual width of a string, excluding ANSI escape sequences
func VisualWidth(s string) int {
	// Remove ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleaned := ansiRegex.ReplaceAllString(s, "")
	return utf8.RuneCountInString(cleaned)
}

// TruncateVisual truncates a string to a specific visual width, preserving colors
func TruncateVisual(s string, width int) string {
	if VisualWidth(s) <= width {
		return s
	}
	
	// This is a simplified version - for full implementation we'd need to
	// properly handle ANSI codes during truncation
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleaned := ansiRegex.ReplaceAllString(s, "")
	runes := []rune(cleaned)
	
	if len(runes) > width {
		return string(runes[:width])
	}
	
	return string(runes)
}