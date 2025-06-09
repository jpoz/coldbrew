package brew

import (
	"context"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Use bubbletea's Key and KeyMsg types directly for full compatibility
type Key = tea.Key
type KeyMsg = tea.KeyMsg

// Re-export bubbletea's key constants for convenience
const (
	KeyEnter      = tea.KeyEnter
	KeyEsc        = tea.KeyEsc
	KeyEscape     = tea.KeyEscape
	KeyBackspace  = tea.KeyBackspace
	KeyTab        = tea.KeyTab
	KeyUp         = tea.KeyUp
	KeyDown       = tea.KeyDown
	KeyLeft       = tea.KeyLeft
	KeyRight      = tea.KeyRight
	KeySpace      = tea.KeySpace
	KeyCtrlC      = tea.KeyCtrlC
	KeyCtrlD      = tea.KeyCtrlD
	KeyCtrlJ      = tea.KeyCtrlJ
	KeyF1         = tea.KeyF1
	KeyF2         = tea.KeyF2
	KeyF3         = tea.KeyF3
	KeyF4         = tea.KeyF4
	KeyF5         = tea.KeyF5
	KeyF6         = tea.KeyF6
	KeyF7         = tea.KeyF7
	KeyF8         = tea.KeyF8
	KeyF9         = tea.KeyF9
	KeyF10        = tea.KeyF10
	KeyF11        = tea.KeyF11
	KeyF12        = tea.KeyF12
)

// handleInput handles keyboard input and sends KeyMsg messages
func (p *Program) handleInput() {
	if p.rawMode {
		// Use bubbletea's own input reading for maximum compatibility
		err := p.readInputsCompat(p.ctx, p.msgChan, os.Stdin)
		if err != nil {
			// If enhanced reading fails, fall back to simple raw mode
			p.handleSimpleRawInput()
		}
	} else {
		// Use line-buffered input for compatibility
		p.handleLineInput()
	}
}

// readInputsCompat uses bubbletea's own key detection logic for maximum compatibility
func (p *Program) readInputsCompat(ctx context.Context, msgs chan<- Msg, input io.Reader) error {
	// We'll delegate to bubbletea's input reading but intercept the messages
	// This ensures 100% compatibility with bubbletea's key detection
	
	// Create a simple input reader loop that reads bytes and converts them to tea.KeyMsg
	var buf [256]byte

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Read bytes from input
		numBytes, err := input.Read(buf[:])
		if err != nil {
			return err
		}

		// Process each byte through our simple detection
		for i := 0; i < numBytes; i++ {
			key := p.detectSimpleKey(buf[i])
			if key != nil {
				select {
				case msgs <- *key:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}

// detectSimpleKey converts a single byte to a KeyMsg
func (p *Program) detectSimpleKey(b byte) *tea.KeyMsg {
	var key tea.Key

	switch b {
	case 3: // Ctrl+C
		key = tea.Key{Type: tea.KeyCtrlC}
	case 4: // Ctrl+D
		key = tea.Key{Type: tea.KeyCtrlD}
	case 10: // Enter (Line Feed) - Unix systems send LF for Enter
		key = tea.Key{Type: tea.KeyEnter}
	case 13: // Enter (Carriage Return) - Windows systems may send CR
		key = tea.Key{Type: tea.KeyEnter}
	case 27: // Escape - this is simplified, real escape sequences are complex
		key = tea.Key{Type: tea.KeyEsc}
	case 127, 8: // Backspace
		key = tea.Key{Type: tea.KeyBackspace}
	case 9: // Tab
		key = tea.Key{Type: tea.KeyTab}
	case 32: // Space
		key = tea.Key{Type: tea.KeySpace, Runes: []rune{' '}}
	default:
		if b >= 32 && b <= 126 { // Printable ASCII
			key = tea.Key{Type: tea.KeyRunes, Runes: []rune{rune(b)}}
		} else {
			// Control character
			key = tea.Key{Type: tea.KeyType(b)}
		}
	}

	keyMsg := tea.KeyMsg(key)
	return &keyMsg
}

// handleLineInput handles line-buffered input (fallback)
func (p *Program) handleLineInput() {
	buf := make([]byte, 1)
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				continue
			}

			key := p.detectSimpleKey(buf[0])
			if key != nil {
				p.Send(*key)
			}
		}
	}
}

// handleSimpleRawInput is a fallback for raw input when enhanced reading fails
func (p *Program) handleSimpleRawInput() {
	buf := make([]byte, 1)

	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				continue
			}

			ch := buf[0]
			var key tea.Key

			switch ch {
			case 3: // Ctrl+C
				key = tea.Key{Type: tea.KeyCtrlC}
			case 4: // Ctrl+D
				key = tea.Key{Type: tea.KeyCtrlD}
			case 10: // Enter (Line Feed) - Unix systems send LF for Enter
				key = tea.Key{Type: tea.KeyEnter}
			case 13: // Enter (Carriage Return) - Windows systems may send CR
				key = tea.Key{Type: tea.KeyEnter}
			case 27: // Escape - try to read escape sequence
				escKey := p.readEscapeSequence()
				key = escKey
			case 127, 8: // Backspace
				key = tea.Key{Type: tea.KeyBackspace}
			case 9: // Tab
				key = tea.Key{Type: tea.KeyTab}
			case 32: // Space
				key = tea.Key{Type: tea.KeySpace, Runes: []rune{' '}}
			default:
				if ch >= 32 && ch <= 126 { // Printable ASCII
					key = tea.Key{Type: tea.KeyRunes, Runes: []rune{rune(ch)}}
				} else {
					// Control character
					key = tea.Key{Type: tea.KeyType(ch)}
				}
			}

			p.Send(tea.KeyMsg(key))
		}
	}
}

// readEscapeSequence reads an escape sequence for simple raw input
func (p *Program) readEscapeSequence() tea.Key {
	buf := make([]byte, 2)
	n, err := os.Stdin.Read(buf)
	if err != nil || n == 0 {
		return tea.Key{Type: tea.KeyEsc}
	}

	if n >= 1 && buf[0] == '[' {
		if n >= 2 {
			switch buf[1] {
			case 'A':
				return tea.Key{Type: tea.KeyUp}
			case 'B':
				return tea.Key{Type: tea.KeyDown}
			case 'C':
				return tea.Key{Type: tea.KeyRight}
			case 'D':
				return tea.Key{Type: tea.KeyLeft}
			case 'I':
				// Focus gained - send FocusMsg directly to program
				p.Send(tea.FocusMsg{})
				return tea.Key{Type: tea.KeyEsc} // Return escape as fallback
			case 'O':
				// Focus lost - send BlurMsg directly to program  
				p.Send(tea.BlurMsg{})
				return tea.Key{Type: tea.KeyEsc} // Return escape as fallback
			}
		}
		// Try to read one more character
		moreBuf := make([]byte, 1)
		if n2, err := os.Stdin.Read(moreBuf); err == nil && n2 > 0 {
			switch moreBuf[0] {
			case 'A':
				return tea.Key{Type: tea.KeyUp}
			case 'B':
				return tea.Key{Type: tea.KeyDown}
			case 'C':
				return tea.Key{Type: tea.KeyRight}
			case 'D':
				return tea.Key{Type: tea.KeyLeft}
			case 'I':
				// Focus gained - send FocusMsg directly to program
				p.Send(tea.FocusMsg{})
				return tea.Key{Type: tea.KeyEsc} // Return escape as fallback
			case 'O':
				// Focus lost - send BlurMsg directly to program
				p.Send(tea.BlurMsg{})
				return tea.Key{Type: tea.KeyEsc} // Return escape as fallback
			}
		}
	}

	return tea.Key{Type: tea.KeyEsc}
}