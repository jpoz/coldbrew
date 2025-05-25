package trmnl

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"unicode/utf8"
)

// KeyType represents the type of key pressed
type KeyType int

// Control keys (reusing bubbletea's key type values for compatibility)
const (
	keyNUL KeyType = 0
	keySOH KeyType = 1
	keySTX KeyType = 2
	keyETX KeyType = 3
	keyEOT KeyType = 4
	keyENQ KeyType = 5
	keyACK KeyType = 6
	keyBEL KeyType = 7
	keyBS  KeyType = 8
	keyHT  KeyType = 9
	keyLF  KeyType = 10
	keyVT  KeyType = 11
	keyFF  KeyType = 12
	keyCR  KeyType = 13
	keySO  KeyType = 14
	keySI  KeyType = 15
	keyDLE KeyType = 16
	keyDC1 KeyType = 17
	keyDC2 KeyType = 18
	keyDC3 KeyType = 19
	keyDC4 KeyType = 20
	keyNAK KeyType = 21
	keySYN KeyType = 22
	keyETB KeyType = 23
	keyCAN KeyType = 24
	keyEM  KeyType = 25
	keySUB KeyType = 26
	keyESC KeyType = 27
	keyFS  KeyType = 28
	keyGS  KeyType = 29
	keyRS  KeyType = 30
	keyUS  KeyType = 31
	keyDEL KeyType = 127
)

// Key type aliases for common keys
const (
	KeyNull      KeyType = keyNUL
	KeyBreak     KeyType = keyETX
	KeyEnter     KeyType = keyCR
	KeyBackspace KeyType = keyDEL
	KeyTab       KeyType = keyHT
	KeyEsc       KeyType = keyESC
	KeyEscape    KeyType = keyESC

	KeyCtrlAt           KeyType = keyNUL
	KeyCtrlA            KeyType = keySOH
	KeyCtrlB            KeyType = keySTX
	KeyCtrlC            KeyType = keyETX
	KeyCtrlD            KeyType = keyEOT
	KeyCtrlE            KeyType = keyENQ
	KeyCtrlF            KeyType = keyACK
	KeyCtrlG            KeyType = keyBEL
	KeyCtrlH            KeyType = keyBS
	KeyCtrlI            KeyType = keyHT
	KeyCtrlJ            KeyType = keyLF
	KeyCtrlK            KeyType = keyVT
	KeyCtrlL            KeyType = keyFF
	KeyCtrlM            KeyType = keyCR
	KeyCtrlN            KeyType = keySO
	KeyCtrlO            KeyType = keySI
	KeyCtrlP            KeyType = keyDLE
	KeyCtrlQ            KeyType = keyDC1
	KeyCtrlR            KeyType = keyDC2
	KeyCtrlS            KeyType = keyDC3
	KeyCtrlT            KeyType = keyDC4
	KeyCtrlU            KeyType = keyNAK
	KeyCtrlV            KeyType = keySYN
	KeyCtrlW            KeyType = keyETB
	KeyCtrlX            KeyType = keyCAN
	KeyCtrlY            KeyType = keyEM
	KeyCtrlZ            KeyType = keySUB
	KeyCtrlOpenBracket  KeyType = keyESC
	KeyCtrlBackslash    KeyType = keyFS
	KeyCtrlCloseBracket KeyType = keyGS
	KeyCtrlCaret        KeyType = keyRS
	KeyCtrlUnderscore   KeyType = keyUS
	KeyCtrlQuestionMark KeyType = keyDEL
)

// Other keys
const (
	KeyRunes KeyType = -(iota + 1)
	KeyUp
	KeyDown
	KeyRight
	KeyLeft
	KeyShiftTab
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDown
	KeyCtrlPgUp
	KeyCtrlPgDown
	KeyDelete
	KeyInsert
	KeySpace
	KeyCtrlUp
	KeyCtrlDown
	KeyCtrlRight
	KeyCtrlLeft
	KeyCtrlHome
	KeyCtrlEnd
	KeyShiftUp
	KeyShiftDown
	KeyShiftRight
	KeyShiftLeft
	KeyShiftHome
	KeyShiftEnd
	KeyCtrlShiftUp
	KeyCtrlShiftDown
	KeyCtrlShiftLeft
	KeyCtrlShiftRight
	KeyCtrlShiftHome
	KeyCtrlShiftEnd
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
)

// Key contains information about a keypress, compatible with bubbletea's Key
type Key struct {
	Type  KeyType
	Runes []rune
	Alt   bool
	Paste bool
}

// KeyMsg represents a keyboard input message, compatible with bubbletea's KeyMsg
type KeyMsg Key

// String returns a string representation for a key message
func (k KeyMsg) String() string {
	return Key(k).String()
}

// String returns a friendly string representation for a key
func (k Key) String() string {
	var result string
	if k.Alt {
		result += "alt+"
	}
	if k.Type == KeyRunes {
		if k.Paste {
			result += "["
		}
		result += string(k.Runes)
		if k.Paste {
			result += "]"
		}
		return result
	} else if s, ok := keyNames[k.Type]; ok {
		return result + s
	}
	return ""
}

// Mappings for control keys and other special keys
var keyNames = map[KeyType]string{
	// Control keys
	keyNUL: "ctrl+@",
	keySOH: "ctrl+a",
	keySTX: "ctrl+b",
	keyETX: "ctrl+c",
	keyEOT: "ctrl+d",
	keyENQ: "ctrl+e",
	keyACK: "ctrl+f",
	keyBEL: "ctrl+g",
	keyBS:  "ctrl+h",
	keyHT:  "tab",
	keyLF:  "ctrl+j",
	keyVT:  "ctrl+k",
	keyFF:  "ctrl+l",
	keyCR:  "enter",
	keySO:  "ctrl+n",
	keySI:  "ctrl+o",
	keyDLE: "ctrl+p",
	keyDC1: "ctrl+q",
	keyDC2: "ctrl+r",
	keyDC3: "ctrl+s",
	keyDC4: "ctrl+t",
	keyNAK: "ctrl+u",
	keySYN: "ctrl+v",
	keyETB: "ctrl+w",
	keyCAN: "ctrl+x",
	keyEM:  "ctrl+y",
	keySUB: "ctrl+z",
	keyESC: "esc",
	keyFS:  "ctrl+\\",
	keyGS:  "ctrl+]",
	keyRS:  "ctrl+^",
	keyUS:  "ctrl+_",
	keyDEL: "backspace",

	// Other keys
	KeyRunes:          "runes",
	KeyUp:             "up",
	KeyDown:           "down",
	KeyRight:          "right",
	KeySpace:          " ",
	KeyLeft:           "left",
	KeyShiftTab:       "shift+tab",
	KeyHome:           "home",
	KeyEnd:            "end",
	KeyCtrlHome:       "ctrl+home",
	KeyCtrlEnd:        "ctrl+end",
	KeyShiftHome:      "shift+home",
	KeyShiftEnd:       "shift+end",
	KeyCtrlShiftHome:  "ctrl+shift+home",
	KeyCtrlShiftEnd:   "ctrl+shift+end",
	KeyPgUp:           "pgup",
	KeyPgDown:         "pgdown",
	KeyCtrlPgUp:       "ctrl+pgup",
	KeyCtrlPgDown:     "ctrl+pgdown",
	KeyDelete:         "delete",
	KeyInsert:         "insert",
	KeyCtrlUp:         "ctrl+up",
	KeyCtrlDown:       "ctrl+down",
	KeyCtrlRight:      "ctrl+right",
	KeyCtrlLeft:       "ctrl+left",
	KeyShiftUp:        "shift+up",
	KeyShiftDown:      "shift+down",
	KeyShiftRight:     "shift+right",
	KeyShiftLeft:      "shift+left",
	KeyCtrlShiftUp:    "ctrl+shift+up",
	KeyCtrlShiftDown:  "ctrl+shift+down",
	KeyCtrlShiftLeft:  "ctrl+shift+left",
	KeyCtrlShiftRight: "ctrl+shift+right",
	KeyF1:             "f1",
	KeyF2:             "f2",
	KeyF3:             "f3",
	KeyF4:             "f4",
	KeyF5:             "f5",
	KeyF6:             "f6",
	KeyF7:             "f7",
	KeyF8:             "f8",
	KeyF9:             "f9",
	KeyF10:            "f10",
	KeyF11:            "f11",
	KeyF12:            "f12",
	KeyF13:            "f13",
	KeyF14:            "f14",
	KeyF15:            "f15",
	KeyF16:            "f16",
	KeyF17:            "f17",
	KeyF18:            "f18",
	KeyF19:            "f19",
	KeyF20:            "f20",
}

// Sequence mappings (reimplemented from bubbletea)
var sequences = map[string]Key{
	// Arrow keys
	"\x1b[A":    {Type: KeyUp},
	"\x1b[B":    {Type: KeyDown},
	"\x1b[C":    {Type: KeyRight},
	"\x1b[D":    {Type: KeyLeft},
	"\x1b[1;2A": {Type: KeyShiftUp},
	"\x1b[1;2B": {Type: KeyShiftDown},
	"\x1b[1;2C": {Type: KeyShiftRight},
	"\x1b[1;2D": {Type: KeyShiftLeft},
	"\x1b[OA":   {Type: KeyShiftUp},
	"\x1b[OB":   {Type: KeyShiftDown},
	"\x1b[OC":   {Type: KeyShiftRight},
	"\x1b[OD":   {Type: KeyShiftLeft},
	"\x1b[a":    {Type: KeyShiftUp},
	"\x1b[b":    {Type: KeyShiftDown},
	"\x1b[c":    {Type: KeyShiftRight},
	"\x1b[d":    {Type: KeyShiftLeft},
	"\x1b[1;3A": {Type: KeyUp, Alt: true},
	"\x1b[1;3B": {Type: KeyDown, Alt: true},
	"\x1b[1;3C": {Type: KeyRight, Alt: true},
	"\x1b[1;3D": {Type: KeyLeft, Alt: true},

	"\x1b[1;4A": {Type: KeyShiftUp, Alt: true},
	"\x1b[1;4B": {Type: KeyShiftDown, Alt: true},
	"\x1b[1;4C": {Type: KeyShiftRight, Alt: true},
	"\x1b[1;4D": {Type: KeyShiftLeft, Alt: true},

	"\x1b[1;5A": {Type: KeyCtrlUp},
	"\x1b[1;5B": {Type: KeyCtrlDown},
	"\x1b[1;5C": {Type: KeyCtrlRight},
	"\x1b[1;5D": {Type: KeyCtrlLeft},
	"\x1b[Oa":   {Type: KeyCtrlUp, Alt: true},
	"\x1b[Ob":   {Type: KeyCtrlDown, Alt: true},
	"\x1b[Oc":   {Type: KeyCtrlRight, Alt: true},
	"\x1b[Od":   {Type: KeyCtrlLeft, Alt: true},
	"\x1b[1;6A": {Type: KeyCtrlShiftUp},
	"\x1b[1;6B": {Type: KeyCtrlShiftDown},
	"\x1b[1;6C": {Type: KeyCtrlShiftRight},
	"\x1b[1;6D": {Type: KeyCtrlShiftLeft},
	"\x1b[1;7A": {Type: KeyCtrlUp, Alt: true},
	"\x1b[1;7B": {Type: KeyCtrlDown, Alt: true},
	"\x1b[1;7C": {Type: KeyCtrlRight, Alt: true},
	"\x1b[1;7D": {Type: KeyCtrlLeft, Alt: true},
	"\x1b[1;8A": {Type: KeyCtrlShiftUp, Alt: true},
	"\x1b[1;8B": {Type: KeyCtrlShiftDown, Alt: true},
	"\x1b[1;8C": {Type: KeyCtrlShiftRight, Alt: true},
	"\x1b[1;8D": {Type: KeyCtrlShiftLeft, Alt: true},

	// Miscellaneous keys
	"\x1b[Z": {Type: KeyShiftTab},

	"\x1b[2~":   {Type: KeyInsert},
	"\x1b[3;2~": {Type: KeyInsert, Alt: true},

	"\x1b[3~":   {Type: KeyDelete},
	"\x1b[3;3~": {Type: KeyDelete, Alt: true},

	"\x1b[5~":   {Type: KeyPgUp},
	"\x1b[5;3~": {Type: KeyPgUp, Alt: true},
	"\x1b[5;5~": {Type: KeyCtrlPgUp},
	"\x1b[5^": {Type: KeyCtrlPgUp},
	"\x1b[5;7~": {Type: KeyCtrlPgUp, Alt: true},

	"\x1b[6~":   {Type: KeyPgDown},
	"\x1b[6;3~": {Type: KeyPgDown, Alt: true},
	"\x1b[6;5~": {Type: KeyCtrlPgDown},
	"\x1b[6^": {Type: KeyCtrlPgDown},
	"\x1b[6;7~": {Type: KeyCtrlPgDown, Alt: true},

	"\x1b[1~":   {Type: KeyHome},
	"\x1b[H":    {Type: KeyHome},
	"\x1b[1;3H": {Type: KeyHome, Alt: true},
	"\x1b[1;5H": {Type: KeyCtrlHome},
	"\x1b[1;7H": {Type: KeyCtrlHome, Alt: true},
	"\x1b[1;2H": {Type: KeyShiftHome},
	"\x1b[1;4H": {Type: KeyShiftHome, Alt: true},
	"\x1b[1;6H": {Type: KeyCtrlShiftHome},
	"\x1b[1;8H": {Type: KeyCtrlShiftHome, Alt: true},

	"\x1b[4~":   {Type: KeyEnd},
	"\x1b[F":    {Type: KeyEnd},
	"\x1b[1;3F": {Type: KeyEnd, Alt: true},
	"\x1b[1;5F": {Type: KeyCtrlEnd},
	"\x1b[1;7F": {Type: KeyCtrlEnd, Alt: true},
	"\x1b[1;2F": {Type: KeyShiftEnd},
	"\x1b[1;4F": {Type: KeyShiftEnd, Alt: true},
	"\x1b[1;6F": {Type: KeyCtrlShiftEnd},
	"\x1b[1;8F": {Type: KeyCtrlShiftEnd, Alt: true},

	"\x1b[7~": {Type: KeyHome},
	"\x1b[7^": {Type: KeyCtrlHome},
	"\x1b[7$": {Type: KeyShiftHome},
	"\x1b[7@": {Type: KeyCtrlShiftHome},

	"\x1b[8~": {Type: KeyEnd},
	"\x1b[8^": {Type: KeyCtrlEnd},
	"\x1b[8$": {Type: KeyShiftEnd},
	"\x1b[8@": {Type: KeyCtrlShiftEnd},

	// Function keys
	"\x1b[[A": {Type: KeyF1},
	"\x1b[[B": {Type: KeyF2},
	"\x1b[[C": {Type: KeyF3},
	"\x1b[[D": {Type: KeyF4},
	"\x1b[[E": {Type: KeyF5},

	"\x1bOP": {Type: KeyF1},
	"\x1bOQ": {Type: KeyF2},
	"\x1bOR": {Type: KeyF3},
	"\x1bOS": {Type: KeyF4},

	"\x1b[1;3P": {Type: KeyF1, Alt: true},
	"\x1b[1;3Q": {Type: KeyF2, Alt: true},
	"\x1b[1;3R": {Type: KeyF3, Alt: true},
	"\x1b[1;3S": {Type: KeyF4, Alt: true},

	"\x1b[11~": {Type: KeyF1},
	"\x1b[12~": {Type: KeyF2},
	"\x1b[13~": {Type: KeyF3},
	"\x1b[14~": {Type: KeyF4},

	"\x1b[15~": {Type: KeyF5},
	"\x1b[15;3~": {Type: KeyF5, Alt: true},

	"\x1b[17~": {Type: KeyF6},
	"\x1b[18~": {Type: KeyF7},
	"\x1b[19~": {Type: KeyF8},
	"\x1b[20~": {Type: KeyF9},
	"\x1b[21~": {Type: KeyF10},

	"\x1b[17;3~": {Type: KeyF6, Alt: true},
	"\x1b[18;3~": {Type: KeyF7, Alt: true},
	"\x1b[19;3~": {Type: KeyF8, Alt: true},
	"\x1b[20;3~": {Type: KeyF9, Alt: true},
	"\x1b[21;3~": {Type: KeyF10, Alt: true},

	"\x1b[23~": {Type: KeyF11},
	"\x1b[24~": {Type: KeyF12},

	"\x1b[23;3~": {Type: KeyF11, Alt: true},
	"\x1b[24;3~": {Type: KeyF12, Alt: true},

	"\x1b[1;2P": {Type: KeyF13},
	"\x1b[1;2Q": {Type: KeyF14},

	"\x1b[25~": {Type: KeyF13},
	"\x1b[26~": {Type: KeyF14},

	"\x1b[25;3~": {Type: KeyF13, Alt: true},
	"\x1b[26;3~": {Type: KeyF14, Alt: true},

	"\x1b[1;2R": {Type: KeyF15},
	"\x1b[1;2S": {Type: KeyF16},

	"\x1b[28~": {Type: KeyF15},
	"\x1b[29~": {Type: KeyF16},

	"\x1b[28;3~": {Type: KeyF15, Alt: true},
	"\x1b[29;3~": {Type: KeyF16, Alt: true},

	"\x1b[15;2~": {Type: KeyF17},
	"\x1b[17;2~": {Type: KeyF18},
	"\x1b[18;2~": {Type: KeyF19},
	"\x1b[19;2~": {Type: KeyF20},

	"\x1b[31~": {Type: KeyF17},
	"\x1b[32~": {Type: KeyF18},
	"\x1b[33~": {Type: KeyF19},
	"\x1b[34~": {Type: KeyF20},

	// Powershell sequences
	"\x1bOA": {Type: KeyUp, Alt: false},
	"\x1bOB": {Type: KeyDown, Alt: false},
	"\x1bOC": {Type: KeyRight, Alt: false},
	"\x1bOD": {Type: KeyLeft, Alt: false},
}

var spaceRunes = []rune{' '}
var unknownCSIRe = regexp.MustCompile(`^\x1b\[[\x30-\x3f]*[\x20-\x2f]*[\x40-\x7e]`)

// extSequences contains sequences plus their alternatives with escape prefix
var extSequences = func() map[string]Key {
	s := map[string]Key{}
	for seq, key := range sequences {
		key := key
		s[seq] = key
		if !key.Alt {
			key.Alt = true
			s["\x1b"+seq] = key
		}
	}
	for i := keyNUL + 1; i <= keyDEL; i++ {
		if i == keyESC {
			continue
		}
		s[string([]byte{byte(i)})] = Key{Type: i}
		s[string([]byte{'\x1b', byte(i)})] = Key{Type: i, Alt: true}
		if i == keyUS {
			i = keyDEL - 1
		}
	}
	s[" "] = Key{Type: KeySpace, Runes: spaceRunes}
	s["\x1b "] = Key{Type: KeySpace, Alt: true, Runes: spaceRunes}
	s["\x1b\x1b"] = Key{Type: KeyEscape, Alt: true}
	return s
}()

// unknownInputByteMsg is for invalid utf-8 bytes
type unknownInputByteMsg byte

func (u unknownInputByteMsg) String() string {
	return fmt.Sprintf("?%#02x?", int(u))
}

// unknownCSISequenceMsg is for unrecognized CSI sequences
type unknownCSISequenceMsg []byte

func (u unknownCSISequenceMsg) String() string {
	return fmt.Sprintf("?CSI%+v?", []byte(u)[2:])
}

// handleInput handles keyboard input and sends KeyMsg messages
func (p *Program) handleInput() {
	if p.rawMode {
		// Use enhanced input reading for raw mode
		err := p.readAnsiInputs(p.ctx, p.msgChan, os.Stdin)
		if err != nil {
			// If enhanced reading fails, fall back to simple raw mode
			p.handleSimpleRawInput()
		}
	} else {
		// Use line-buffered input for compatibility
		p.handleLineInput()
	}
}

// readAnsiInputs reads keypress inputs from a TTY and produces KeyMsg messages
// This is a reimplementation of bubbletea's readAnsiInputs for trmnl
func (p *Program) readAnsiInputs(ctx context.Context, msgs chan<- Msg, input io.Reader) error {
	var buf [256]byte

	var leftOverFromPrevIteration []byte
loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Read and block
		numBytes, err := input.Read(buf[:])
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}
		b := buf[:numBytes]
		if leftOverFromPrevIteration != nil {
			b = append(leftOverFromPrevIteration, b...)
		}

		canHaveMoreData := numBytes == len(buf)

		var i, w int
		for i, w = 0, 0; i < len(b); i += w {
			var msg Msg
			w, msg = p.detectOneMsg(b[i:], canHaveMoreData)
			if w == 0 {
				leftOverFromPrevIteration = make([]byte, 0, len(b[i:])+len(buf))
				leftOverFromPrevIteration = append(leftOverFromPrevIteration, b[i:]...)
				continue loop
			}

			select {
			case msgs <- msg:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		leftOverFromPrevIteration = nil
	}
}

// detectOneMsg detects a single message from input bytes
func (p *Program) detectOneMsg(b []byte, canHaveMoreData bool) (w int, msg Msg) {
	// Detect escape sequence and control characters
	var foundSeq bool
	foundSeq, w, msg = p.detectSequence(b)
	if foundSeq {
		return w, msg
	}

	// Handle alt modifier
	alt := false
	i := 0
	if b[0] == '\x1b' {
		alt = true
		i++
	}

	// Handle standalone NUL
	if i < len(b) && b[i] == 0 {
		return i + 1, KeyMsg{Type: keyNUL, Alt: alt}
	}

	// Find the longest sequence of runes
	var runes []rune
	for rw := 0; i < len(b); i += rw {
		var r rune
		r, rw = utf8.DecodeRune(b[i:])
		if r == utf8.RuneError || r <= rune(keyUS) || r == rune(keyDEL) || r == ' ' {
			break
		}
		runes = append(runes, r)
		if alt {
			i += rw
			break
		}
	}
	if i >= len(b) && canHaveMoreData {
		return 0, nil
	}

	// If we found at least one rune, report it
	if len(runes) > 0 {
		k := Key{Type: KeyRunes, Runes: runes, Alt: alt}
		if len(runes) == 1 && runes[0] == ' ' {
			k.Type = KeySpace
		}
		return i, KeyMsg(k)
	}

	// Handle lone escape character
	if alt && len(b) == 1 {
		return 1, KeyMsg(Key{Type: KeyEscape})
	}

	// Invalid byte
	return 1, unknownInputByteMsg(b[0])
}

// detectSequence detects escape sequences and control characters
func (p *Program) detectSequence(input []byte) (hasSeq bool, width int, msg Msg) {
	// Check known sequences
	for seq, key := range extSequences {
		if len(input) >= len(seq) && string(input[:len(seq)]) == seq {
			return true, len(seq), KeyMsg(key)
		}
	}

	// Check for unknown CSI sequence
	if loc := unknownCSIRe.FindIndex(input); loc != nil {
		return true, loc[1], unknownCSISequenceMsg(input[:loc[1]])
	}

	return false, 0, nil
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

			ch := buf[0]
			var key Key

			switch ch {
			case 3: // Ctrl+C
				key = Key{Type: keyETX}
			case 4: // Ctrl+D
				key = Key{Type: keyEOT}
			case 13: // Enter
				key = Key{Type: keyCR}
			case 27: // Escape
				key = Key{Type: keyESC}
			case 127, 8: // Backspace
				key = Key{Type: keyDEL}
			case 9: // Tab
				key = Key{Type: keyHT}
			case 32: // Space
				key = Key{Type: KeySpace, Runes: []rune{' '}}
			default:
				if ch >= 32 && ch <= 126 { // Printable ASCII
					key = Key{Type: KeyRunes, Runes: []rune{rune(ch)}}
				} else {
					// Control character
					key = Key{Type: KeyType(ch)}
				}
			}

			p.Send(KeyMsg(key))
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
			var key Key

			switch ch {
			case 3: // Ctrl+C
				key = Key{Type: keyETX}
			case 4: // Ctrl+D
				key = Key{Type: keyEOT}
			case 13: // Enter
				key = Key{Type: keyCR}
			case 27: // Escape - try to read escape sequence
				escKey := p.readEscapeSequence()
				key = escKey
			case 127, 8: // Backspace
				key = Key{Type: keyDEL}
			case 9: // Tab
				key = Key{Type: keyHT}
			case 32: // Space
				key = Key{Type: KeySpace, Runes: []rune{' '}}
			default:
				if ch >= 32 && ch <= 126 { // Printable ASCII
					key = Key{Type: KeyRunes, Runes: []rune{rune(ch)}}
				} else {
					// Control character
					key = Key{Type: KeyType(ch)}
				}
			}

			p.Send(KeyMsg(key))
		}
	}
}

// readEscapeSequence reads an escape sequence for simple raw input
func (p *Program) readEscapeSequence() Key {
	buf := make([]byte, 2)
	n, err := os.Stdin.Read(buf)
	if err != nil || n == 0 {
		return Key{Type: keyESC}
	}

	if n >= 1 && buf[0] == '[' {
		if n >= 2 {
			switch buf[1] {
			case 'A':
				return Key{Type: KeyUp}
			case 'B':
				return Key{Type: KeyDown}
			case 'C':
				return Key{Type: KeyRight}
			case 'D':
				return Key{Type: KeyLeft}
			}
		}
		// Try to read one more character
		moreBuf := make([]byte, 1)
		if n2, err := os.Stdin.Read(moreBuf); err == nil && n2 > 0 {
			switch moreBuf[0] {
			case 'A':
				return Key{Type: KeyUp}
			case 'B':
				return Key{Type: KeyDown}
			case 'C':
				return Key{Type: KeyRight}
			case 'D':
				return Key{Type: KeyLeft}
			}
		}
	}

	return Key{Type: keyESC}
}