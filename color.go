package trmnl

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
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
	
	visualPos := 0
	runeIndex := 0
	
	for runeIndex < len(runes) {
		// Find if current visual position has color
		var currentColor Color = ColorDefault
		for _, colorInfo := range colorMap {
			if visualPos >= colorInfo.Start && visualPos < colorInfo.End {
				currentColor = colorInfo.Color
				break
			}
		}
		
		if currentColor != ColorDefault {
			// Find the end of this color section by visual position
			result.WriteString(string(currentColor))
			colorEnd := -1
			for _, colorInfo := range colorMap {
				if visualPos >= colorInfo.Start && visualPos < colorInfo.End && colorInfo.Color == currentColor {
					colorEnd = colorInfo.End
					break
				}
			}
			
			// Add runes until we reach the color end position
			for runeIndex < len(runes) && visualPos < colorEnd {
				r := runes[runeIndex]
				result.WriteRune(r)
				visualPos += runewidth.RuneWidth(r)
				runeIndex++
			}
			
			result.WriteString(string(ColorReset))
		} else {
			r := runes[runeIndex]
			result.WriteRune(r)
			visualPos += runewidth.RuneWidth(r)
			runeIndex++
		}
	}
	
	return result.String()
}

// VisualWidth calculates the visual width of a string, excluding ANSI escape sequences
func VisualWidth(s string) int {
	// Remove ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleaned := ansiRegex.ReplaceAllString(s, "")
	
	return runewidth.StringWidth(cleaned)
}

// RuneVisualWidth returns the visual width of a single rune
func RuneVisualWidth(r rune) int {
	// Control characters and non-printable characters
	if r < 32 || (r >= 0x7F && r < 0xA0) {
		return 0
	}
	
	// Common emoji ranges (these are typically 2 columns wide)
	if isEmoji(r) {
		return 2
	}
	
	// East Asian Wide characters (2 columns)
	if isEastAsianWide(r) {
		return 2
	}
	
	// Combining characters (0 columns)
	if unicode.In(r, unicode.Mn, unicode.Me, unicode.Mc) {
		return 0
	}
	
	// Default: most characters are 1 column
	return 1
}

// isEmoji checks if a rune is likely an emoji
func isEmoji(r rune) bool {
	// Common emoji Unicode ranges
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x1F1E6 && r <= 0x1F1FF) || // Regional indicators
		(r >= 0x2600 && r <= 0x26FF) ||   // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) ||   // Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) ||   // Variation selectors
		r == 0x200D ||                     // Zero width joiner
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1FA70 && r <= 0x1FAFF)    // Symbols and Pictographs Extended-A
}

// isEastAsianWide checks if a rune is an East Asian wide character
func isEastAsianWide(r rune) bool {
	// Simplified check for East Asian wide characters
	// Full implementation would use unicode/width package
	return (r >= 0x1100 && r <= 0x115F) || // Hangul Jamo
		(r >= 0x2E80 && r <= 0x2EFF) || // CJK Radicals Supplement
		(r >= 0x2F00 && r <= 0x2FDF) || // Kangxi Radicals
		(r >= 0x3000 && r <= 0x303F) || // CJK Symbols and Punctuation
		(r >= 0x3040 && r <= 0x309F) || // Hiragana
		(r >= 0x30A0 && r <= 0x30FF) || // Katakana
		(r >= 0x3100 && r <= 0x312F) || // Bopomofo
		(r >= 0x3130 && r <= 0x318F) || // Hangul Compatibility Jamo
		(r >= 0x3190 && r <= 0x319F) || // Kanbun
		(r >= 0x31A0 && r <= 0x31BF) || // Bopomofo Extended
		(r >= 0x31C0 && r <= 0x31EF) || // CJK Strokes
		(r >= 0x31F0 && r <= 0x31FF) || // Katakana Phonetic Extensions
		(r >= 0x3200 && r <= 0x32FF) || // Enclosed CJK Letters and Months
		(r >= 0x3300 && r <= 0x33FF) || // CJK Compatibility
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Unified Ideographs Extension A
		(r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0xA000 && r <= 0xA48F) || // Yi Syllables
		(r >= 0xA490 && r <= 0xA4CF) || // Yi Radicals
		(r >= 0xAC00 && r <= 0xD7AF) || // Hangul Syllables
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
		(r >= 0xFE30 && r <= 0xFE4F) || // CJK Compatibility Forms
		(r >= 0xFF00 && r <= 0xFFEF)    // Halfwidth and Fullwidth Forms
}

// TruncateVisual truncates a string to a specific visual width, preserving colors
func TruncateVisual(s string, width int) string {
	if VisualWidth(s) <= width {
		return s
	}
	
	// Remove ANSI escape sequences for width calculation
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleaned := ansiRegex.ReplaceAllString(s, "")
	
	return runewidth.Truncate(cleaned, width, "")
}