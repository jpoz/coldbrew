package trmnl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// KeyMsg represents a keyboard input message
type KeyMsg struct {
	Key  string
	Rune rune
}

// String returns the key as a string
func (k KeyMsg) String() string {
	return k.Key
}

// handleInput handles keyboard input and sends KeyMsg messages
func (p *Program) handleInput() {
	if p.rawMode {
		p.handleRawInput()
	} else {
		p.handleLineInput()
	}
}

// handleLineInput handles line-buffered input (original behavior)
func (p *Program) handleLineInput() {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			// Read a line of input
			input, err := reader.ReadString('\n')
			if err != nil {
				continue
			}
			
			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}
			
			// Handle special keys
			var key string
			var r rune
			
			switch input {
			case "q", "quit", "exit":
				key = "q"
			case "up":
				key = "up"
			case "down":
				key = "down"
			case "left":
				key = "left"
			case "right":
				key = "right"
			case "enter", "":
				key = "enter"
			default:
				if len(input) == 1 {
					r = rune(input[0])
					key = input
				} else {
					key = input
				}
			}
			
			p.Send(KeyMsg{Key: key, Rune: r})
		}
	}
}

// handleRawInput handles immediate character input without Enter
func (p *Program) handleRawInput() {
	buf := make([]byte, 1)
	
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			// Read single character
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				continue
			}
			
			ch := buf[0]
			var key string
			var r rune
			
			// Handle special characters
			switch ch {
			case 3: // Ctrl+C
				key = "ctrl+c"
			case 4: // Ctrl+D (EOF)
				key = "ctrl+d"
			case 13: // Enter/Return
				key = "enter"
			case 27: // Escape sequences (arrow keys, etc.)
				key = p.handleEscapeSequence()
			case 127, 8: // Backspace/Delete
				key = "backspace"
			case 9: // Tab
				key = "tab"
			case 32: // Space
				key = "space"
				r = ' '
			default:
				if ch >= 32 && ch <= 126 { // Printable ASCII
					r = rune(ch)
					key = string(r)
				} else {
					key = fmt.Sprintf("unknown_%d", ch)
				}
			}
			
			p.Send(KeyMsg{Key: key, Rune: r})
		}
	}
}

// handleEscapeSequence handles ANSI escape sequences for arrow keys, etc.
func (p *Program) handleEscapeSequence() string {
	buf := make([]byte, 2)
	
	// Try to read the next characters in the sequence
	n, err := os.Stdin.Read(buf)
	if err != nil || n == 0 {
		return "escape"
	}
	
	if n >= 1 && buf[0] == '[' {
		// Arrow keys and other escape sequences
		if n >= 2 {
			switch buf[1] {
			case 'A':
				return "up"
			case 'B':
				return "down"
			case 'C':
				return "right"
			case 'D':
				return "left"
			}
		}
		
		// Try to read one more character for other sequences
		moreBuf := make([]byte, 1)
		if n2, err := os.Stdin.Read(moreBuf); err == nil && n2 > 0 {
			switch moreBuf[0] {
			case 'A':
				return "up"
			case 'B':
				return "down"
			case 'C':
				return "right"
			case 'D':
				return "left"
			}
		}
	}
	
	return "escape"
}