package trmnl

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Terminal provides rendering utilities
type Terminal struct{}

func NewTerminal() *Terminal {
	return &Terminal{}
}

func (t *Terminal) Render(component Component, size Size) {
	lines := component.Render(size)
	for _, line := range lines {
		fmt.Println(line)
	}
}

func (t *Terminal) Clear() {
	fmt.Print("\033[H\033[2J")
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

// RenderResponsive renders the component using the full terminal size
func (t *Terminal) RenderResponsive(component Component) {
	size, err := t.GetSize()
	if err != nil {
		// Fallback to default size
		size = Size{Width: 80, Height: 24}
	}
	t.Render(component, size)
}

// RenderFullWidth renders the component using the full terminal width with specified height
func (t *Terminal) RenderFullWidth(component Component, height int) {
	size, err := t.GetSize()
	if err != nil {
		// Fallback to default width
		size.Width = 80
	}
	size.Height = height
	t.Render(component, size)
}

// RenderFullHeight renders the component using the full terminal height with specified width
func (t *Terminal) RenderFullHeight(component Component, width int) {
	size, err := t.GetSize()
	if err != nil {
		// Fallback to default height
		size.Height = 24
	}
	size.Width = width
	t.Render(component, size)
}

