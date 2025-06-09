package brew

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"
)

// Terminal provides rendering utilities
type Terminal struct {
	previousBuffer []string
	lastSize       Size
	renderStartRow int
	renderStartCol int
	totalRendered  int
	firstRender    bool
}

func NewTerminal() *Terminal {
	return &Terminal{firstRender: true}
}

func (t *Terminal) Clear() {
	fmt.Print("\033[H\033[2J")
}

// HideCursor hides the terminal cursor
func (t *Terminal) HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func (t *Terminal) ShowCursor() {
	fmt.Print("\033[?25h")
}

// EnableReportFocus enables terminal focus reporting
func (t *Terminal) EnableReportFocus() {
	fmt.Print("\033[?1004h")
}

// DisableReportFocus disables terminal focus reporting
func (t *Terminal) DisableReportFocus() {
	fmt.Print("\033[?1004l")
}

// MoveCursor moves the cursor to a specific position (1-based coordinates)
func (t *Terminal) MoveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

// MoveCursorHome moves the cursor to the top-left corner
func (t *Terminal) MoveCursorHome() {
	fmt.Print("\033[H")
}


// GetSize returns the current terminal dimensions
func (t *Terminal) GetSize() (Size, error) {
	// Try to get size using term.GetSize (same method as bubbletea)
	width, height, err := term.GetSize(os.Stdout.Fd())
	if err == nil {
		return Size{Width: width, Height: height}, nil
	}

	// Fallback to stty if term.GetSize fails
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		// Fallback to default size if detection fails
		return Size{Width: 80, Height: 24}, err
	}

	parts := strings.Fields(strings.TrimSpace(string(output)))
	if len(parts) != 2 {
		return Size{Width: 80, Height: 24}, fmt.Errorf("unexpected stty output format")
	}

	sttyHeight, err := strconv.Atoi(parts[0])
	if err != nil {
		return Size{Width: 80, Height: 24}, err
	}

	sttyWidth, err := strconv.Atoi(parts[1])
	if err != nil {
		return Size{Width: 80, Height: 24}, err
	}

	return Size{Width: sttyWidth, Height: sttyHeight}, nil
}

// RenderString renders a string directly to the terminal with differential updates
func (t *Terminal) RenderString(content string) {
	lines := strings.Split(content, "\n")
	
	// Check for size changes to force full re-render
	currentSize, _ := t.GetSize()
	sizeChanged := t.lastSize.Width != currentSize.Width || t.lastSize.Height != currentSize.Height
	t.lastSize = currentSize
	
	// First render - just render content from current cursor position
	if t.firstRender {
		t.firstRender = false
		
		// Render all content without clearing screen
		for i, line := range lines {
			if i > 0 {
				fmt.Print("\n")
			}
			fmt.Print(line)
		}
		
		// Update state
		t.previousBuffer = make([]string, len(lines))
		copy(t.previousBuffer, lines)
		t.totalRendered = len(lines)
		return
	}
	
	// Size changed - clear screen and re-render to avoid artifacts
	if sizeChanged {
		fmt.Print("\033[H\033[2J") // Clear screen and go to home
		
		// Render all content
		for i, line := range lines {
			if i > 0 {
				fmt.Print("\n")
			}
			fmt.Print(line)
		}
		
		// Update state
		t.previousBuffer = make([]string, len(lines))
		copy(t.previousBuffer, lines)
		t.totalRendered = len(lines)
		return
	}
	
	// Find first differing line
	firstDiff := -1
	minLen := len(lines)
	if len(t.previousBuffer) < minLen {
		minLen = len(t.previousBuffer)
	}
	
	for i := 0; i < minLen; i++ {
		if lines[i] != t.previousBuffer[i] {
			firstDiff = i
			break
		}
	}
	
	// If lengths differ but common lines are same, start diff at end of common
	if firstDiff == -1 && len(lines) != len(t.previousBuffer) {
		firstDiff = minLen
	}
	
	// No changes needed
	if firstDiff == -1 {
		return
	}
	
	// Move cursor to the first line that needs updating
	// We need to go back up from where we currently are (end of previous render)
	// to the first differing line
	linesToGoBack := t.totalRendered - 1 - firstDiff
	if linesToGoBack > 0 {
		fmt.Printf("\033[%dA", linesToGoBack) // Move up
	}
	fmt.Print("\r") // Move to beginning of line
	
	// Clear from current position to end of screen
	fmt.Print("\033[J")
	
	// Render changed lines
	for i := firstDiff; i < len(lines); i++ {
		if i > firstDiff {
			fmt.Print("\n")
		}
		fmt.Print(lines[i])
	}
	
	// Update state
	t.previousBuffer = make([]string, len(lines))
	copy(t.previousBuffer, lines)
	t.totalRendered = len(lines)
}

// ClearPreviousBuffer clears the stored previous buffer (useful for manual redraws)
func (t *Terminal) ClearPreviousBuffer() {
	t.previousBuffer = nil
}
