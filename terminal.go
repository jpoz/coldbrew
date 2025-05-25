package trmnl

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

// RenderString renders a string directly to the terminal
func (t *Terminal) RenderString(content string) {
	lines := strings.Split(content, "\n")
	t.renderWithDiff(lines)
}

// renderWithDiff performs diff-based rendering to minimize screen updates
func (t *Terminal) renderWithDiff(newLines []string) {
	// Diff-based rendering: only update changed lines
	maxLines := max(len(newLines), len(t.previousBuffer))

	for i := range maxLines {
		var newLine, oldLine string

		if i < len(newLines) {
			newLine = newLines[i]
		}
		if i < len(t.previousBuffer) {
			oldLine = t.previousBuffer[i]
		}

		// Only update if line changed
		if newLine != oldLine {
			t.MoveCursor(i+1, 1) // Move to beginning of line (1-based coordinates)

			if newLine == "" {
				// Clear the entire line
				fmt.Print("\033[K")
			} else {
				// Print new content and clear rest of line
				fmt.Print(newLine)
				fmt.Print("\033[K")
			}
		}
	}

	// Update previous buffer
	t.previousBuffer = make([]string, len(newLines))
	copy(t.previousBuffer, newLines)
}

// ClearPreviousBuffer clears the stored previous buffer (useful for manual redraws)
func (t *Terminal) ClearPreviousBuffer() {
	t.previousBuffer = nil
}
