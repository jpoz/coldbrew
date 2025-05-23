package trmnl

import "fmt"

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