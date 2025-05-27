package brew

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Terminal provides rendering utilities
type Terminal struct {
	previousBuffer []string
	lastSize       Size
	altScreen      bool
}

func NewTerminal() *Terminal {
	return &Terminal{}
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

// MoveCursor moves the cursor to a specific position (1-based coordinates)
func (t *Terminal) MoveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

// MoveCursorHome moves the cursor to the top-left corner
func (t *Terminal) MoveCursorHome() {
	fmt.Print("\033[H")
}

// EnterAltScreen enters the alternate screen buffer
func (t *Terminal) EnterAltScreen() {
	if !t.altScreen {
		fmt.Print("\033[?1049h") // Enter alternate screen
		t.altScreen = true
	}
}

// ExitAltScreen exits the alternate screen buffer
func (t *Terminal) ExitAltScreen() {
	if t.altScreen {
		fmt.Print("\033[?1049l") // Exit alternate screen
		t.altScreen = false
	}
}

// GetSize returns the current terminal dimensions
func (t *Terminal) GetSize() (Size, error) {
	// Try to get size using stty
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

	height, err := strconv.Atoi(parts[0])
	if err != nil {
		return Size{Width: 80, Height: 24}, err
	}

	width, err := strconv.Atoi(parts[1])
	if err != nil {
		return Size{Width: 80, Height: 24}, err
	}

	return Size{Width: width, Height: height}, nil
}

// RenderString renders a string directly to the terminal using alternate screen buffer
func (t *Terminal) RenderString(content string) {
	// Enter alternate screen on first render
	if len(t.previousBuffer) == 0 {
		t.EnterAltScreen()
	}
	
	// Clear screen and move cursor to home position
	t.Clear()
	t.MoveCursorHome()
	
	// Simply print the content - no complex cursor tracking needed in alt screen
	fmt.Print(content)
	
	lines := strings.Split(content, "\n")
	t.previousBuffer = make([]string, len(lines))
	copy(t.previousBuffer, lines)
}

// ClearPreviousBuffer clears the stored previous buffer (useful for manual redraws)
func (t *Terminal) ClearPreviousBuffer() {
	t.previousBuffer = nil
}
