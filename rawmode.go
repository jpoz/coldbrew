package trmnl

import (
	"os"
	"syscall"
	"unsafe"
)

// TerminalState holds the original terminal state
type TerminalState struct {
	fd       int
	original *syscall.Termios
}

// enableRawMode puts the terminal into raw mode for immediate input
func enableRawMode() (*TerminalState, error) {
	fd := int(os.Stdin.Fd())
	
	// Get original terminal attributes
	original, err := getTerminalState(fd)
	if err != nil {
		return nil, err
	}
	
	// Create a copy to modify
	raw := *original
	
	// Disable canonical mode (line buffering) and echo
	raw.Lflag &^= syscall.ICANON | syscall.ECHO
	
	// Set minimum characters to read and timeout
	raw.Cc[syscall.VMIN] = 1  // Read at least 1 character
	raw.Cc[syscall.VTIME] = 0 // No timeout
	
	// Apply the new settings
	if err := setTerminalState(fd, &raw); err != nil {
		return nil, err
	}
	
	return &TerminalState{fd: fd, original: original}, nil
}

// restoreTerminalState restores the terminal to its original state
func (ts *TerminalState) restore() error {
	if ts == nil || ts.original == nil {
		return nil
	}
	return setTerminalState(ts.fd, ts.original)
}

// Platform-specific constants for macOS/Darwin
const (
	TCGETS = 0x40487413
	TCSETS = 0x80487414
)

// getTerminalState gets the current terminal attributes
func getTerminalState(fd int) (*syscall.Termios, error) {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TCGETS, uintptr(unsafe.Pointer(&termios)))
	if errno != 0 {
		return nil, errno
	}
	return &termios, nil
}

// setTerminalState sets terminal attributes
func setTerminalState(fd int, termios *syscall.Termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TCSETS, uintptr(unsafe.Pointer(termios)))
	if errno != 0 {
		return errno
	}
	return nil
}