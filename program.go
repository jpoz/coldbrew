package brew

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Type aliases for bubbletea compatibility
type Msg = tea.Msg
type Cmd = tea.Cmd
type Model = tea.Model

// Sub represents a subscription that can send messages continuously
type Sub func(ctx context.Context, send func(Msg))

// ColdbrewModel extends tea.Model with subscriptions for advanced functionality
type ColdbrewModel interface {
	tea.Model

	// Subscriptions returns active subscriptions (optional)
	Subscriptions() []Sub
}

// Program manages the Elm architecture runtime
type Program struct {
	model         tea.Model
	terminal      *Terminal
	msgChan       chan Msg
	quit          chan struct{}
	ctx           context.Context
	cancel        context.CancelFunc
	subs          []Sub
	hideCursor    bool
	rawMode       bool
	terminalState *TerminalState
	finished      chan struct{}
}

// QuitMsg signals the program should exit
type QuitMsg struct{}

// NewProgram creates a new Elm architecture program
func NewProgram(initialModel tea.Model) *Program {
	ctx, cancel := context.WithCancel(context.Background())
	return &Program{
		model:      initialModel,
		terminal:   NewTerminal(),
		msgChan:    make(chan Msg, 100),
		quit:       make(chan struct{}),
		ctx:        ctx,
		cancel:     cancel,
		hideCursor: true, // Default to hiding cursor for TUI apps
		rawMode:    true, // Default to raw mode for immediate input
		finished:   make(chan struct{}),
	}
}

// WithCursorHidden sets whether the cursor should be hidden during program execution
func (p *Program) WithCursorHidden(hide bool) *Program {
	p.hideCursor = hide
	return p
}

// WithRawMode enables immediate character input without pressing Enter
func (p *Program) WithRawMode(enable bool) *Program {
	p.rawMode = enable
	return p
}

// Send sends a message to the program
func (p *Program) Send(msg Msg) {
	select {
	case p.msgChan <- msg:
	case <-p.ctx.Done():
	}
}

// Quit is a convenience function for quitting Bubble Tea programs. Use it
// when you need to shut down a Bubble Tea program from the outside.
//
// If you wish to quit from within a Bubble Tea program use the Quit command.
//
// If the program is not running this will be a no-op, so it's safe to call
// if the program is unstarted or has already exited.
// This is compatible with bubbletea's Program.Quit method.
func (p *Program) Quit() {
	p.Send(tea.Quit())
}

// Kill signals the program to stop immediately and restore the former terminal state.
// The final render that you would normally see when quitting will be skipped.
// This is compatible with bubbletea's Program.Kill method.
func (p *Program) Kill() {
	p.cancel()
}

// Wait waits/blocks until the underlying Program finished shutting down.
// This is compatible with bubbletea's Program.Wait method.
func (p *Program) Wait() {
	<-p.finished
}

// Run starts the program and blocks until it exits
// Returns the final model and any error, matching bubbletea's signature
func (p *Program) Run() (tea.Model, error) {
	// Ensure finished channel is closed when Run exits
	defer close(p.finished)
	
	// Setup cursor visibility
	if p.hideCursor {
		p.terminal.HideCursor()
		// Ensure cursor is restored on exit
		defer func() {
			p.terminal.ShowCursor()
			// Print a final newline to position cursor properly for shell prompt
			fmt.Print("\n")
		}()
	}

	// Setup raw mode if enabled
	if p.rawMode {
		termState, err := enableRawMode()
		if err != nil {
			// If raw mode fails (e.g., not a real terminal), fall back to line mode
			p.rawMode = false
		} else {
			p.terminalState = termState
			// Ensure terminal state is restored on exit
			defer p.terminalState.restore()
		}
	}

	// Send initial window size
	go p.checkResize()

	// Start listening for window resize events
	go p.handleResize()

	// Initialize the model
	cmd := p.model.Init()

	// Execute initial command if any
	if cmd != nil {
		go func() {
			if msg := cmd(); msg != nil {
				p.Send(msg)
			}
		}()
	}

	// Initial render
	p.render()

	// Start input handling
	go p.handleInput()

	// Start subscriptions
	p.startSubscriptions()

	// Handle interrupt signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		p.Send(tea.Quit())
	}()

	// Main message loop
	for {
		select {
		case msg := <-p.msgChan:
			// Handle quit messages - both tea.Quit and custom QuitMsg
			if _, isQuit := msg.(QuitMsg); isQuit {
				p.cancel()
				return p.model, nil
			}
			if _, isTeaQuit := msg.(tea.QuitMsg); isTeaQuit {
				p.cancel()
				return p.model, nil
			}

			// Handle windowSizeMsg (internal message to trigger size check)
			if _, isWindowSizeMsg := msg.(windowSizeMsg); isWindowSizeMsg {
				go p.checkResize()
				continue
			}

			// Handle enableReportFocusMsg (enable focus reporting)
			if _, isEnableFocus := msg.(enableReportFocusMsg); isEnableFocus {
				p.terminal.EnableReportFocus()
				continue
			}

			// Handle bubbletea's enableReportFocusMsg (check by comparing with tea.EnableReportFocus())
			if fmt.Sprintf("%T", msg) == fmt.Sprintf("%T", tea.EnableReportFocus()) {
				p.terminal.EnableReportFocus()
				continue
			}

			// Handle disableReportFocusMsg (disable focus reporting)
			if _, isDisableFocus := msg.(disableReportFocusMsg); isDisableFocus {
				p.terminal.DisableReportFocus()
				continue
			}

			// Handle bubbletea's disableReportFocusMsg (check by comparing with tea.DisableReportFocus())
			if fmt.Sprintf("%T", msg) == fmt.Sprintf("%T", tea.DisableReportFocus()) {
				p.terminal.DisableReportFocus()
				continue
			}

			// Handle tea.BatchMsg (bubbletea's batch commands)
			if teaBatchMsg, isTeaBatch := msg.(tea.BatchMsg); isTeaBatch {
				for _, cmd := range teaBatchMsg {
					if cmd != nil {
						go func(cmd tea.Cmd) {
							if cmdMsg := cmd(); cmdMsg != nil {
								p.Send(cmdMsg)
							}
						}(cmd)
					}
				}
				continue
			}

			// Handle coldbrew BatchMsg (legacy)
			if batchMsg, isBatch := msg.(BatchMsg); isBatch {
				for _, batchedMsg := range batchMsg.Messages {
					// Process each batched message synchronously
					newModel, newCmd := p.model.Update(batchedMsg)
					p.model = newModel

					// Execute command if any
					if newCmd != nil {
						go func() {
							if cmdMsg := newCmd(); cmdMsg != nil {
								p.Send(cmdMsg)
							}
						}()
					}
				}
				// Re-render after all batch messages processed
				p.render()
				continue
			}

			// Update model
			newModel, newCmd := p.model.Update(msg)
			p.model = newModel

			// Execute command if any
			if newCmd != nil {
				go func() {
					if cmdMsg := newCmd(); cmdMsg != nil {
						p.Send(cmdMsg)
					}
				}()
			}

			// Re-render
			p.render()

		case <-p.ctx.Done():
			return p.model, nil
		}
	}
}

// render renders the current model to the terminal
func (p *Program) render() {
	viewString := p.model.View()
	p.terminal.RenderString(viewString)
}

// Quit creates a command that quits the program
func Quit() Cmd {
	return func() Msg {
		return QuitMsg{}
	}
}

// BatchMsg represents multiple messages to be processed
type BatchMsg struct {
	Messages []Msg
}

// Batch creates a command that sends multiple messages
func Batch(messages ...Msg) Cmd {
	return func() Msg {
		return BatchMsg{Messages: messages}
	}
}

// Delay creates a command that sends a message after a delay
func Delay(duration time.Duration, msg Msg) Cmd {
	return func() Msg {
		time.Sleep(duration)
		return msg
	}
}

// Tick creates a command that sends a message after a duration
func Tick(duration time.Duration, msgFunc func(time.Time) Msg) Cmd {
	return func() Msg {
		time.Sleep(duration)
		return msgFunc(time.Now())
	}
}

// WindowSize creates a command that queries the terminal for its current size
// This is compatible with bubbletea's WindowSize command
func WindowSize() Cmd {
	return func() Msg {
		return windowSizeMsg{}
	}
}

// windowSizeMsg is used internally to trigger a window size check
type windowSizeMsg struct{}

// EnableReportFocus enables terminal focus reporting
// This is compatible with bubbletea's EnableReportFocus command
func EnableReportFocus() Cmd {
	return func() Msg {
		return enableReportFocusMsg{}
	}
}

// enableReportFocusMsg is used internally to enable focus reporting
type enableReportFocusMsg struct{}

// DisableReportFocus disables terminal focus reporting
// This is compatible with bubbletea's DisableReportFocus command
func DisableReportFocus() Cmd {
	return func() Msg {
		return disableReportFocusMsg{}
	}
}

// disableReportFocusMsg is used internally to disable focus reporting
type disableReportFocusMsg struct{}

// startSubscriptions starts all subscriptions from the model if it supports them
func (p *Program) startSubscriptions() {
	// Check if model supports subscriptions
	if coldbrewModel, ok := p.model.(ColdbrewModel); ok {
		p.subs = coldbrewModel.Subscriptions()

		// Start each subscription in its own goroutine
		for _, sub := range p.subs {
			go sub(p.ctx, p.Send)
		}
	}
}

// Every creates a subscription that sends messages at regular intervals
func Every(duration time.Duration, msgFunc func(time.Time) Msg) Sub {
	return func(ctx context.Context, send func(Msg)) {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				send(msgFunc(t))
			case <-ctx.Done():
				return
			}
		}
	}
}

// handleResize listens for terminal resize events and sends WindowSizeMsg
func (p *Program) handleResize() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	defer signal.Stop(sig)

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-sig:
			p.checkResize()
		}
	}
}

// checkResize detects the current terminal size and sends a WindowSizeMsg
func (p *Program) checkResize() {
	size, err := p.terminal.GetSize()
	if err != nil {
		// Silently ignore errors (terminal may not support size detection)
		return
	}

	p.Send(tea.WindowSizeMsg{
		Width:  size.Width,
		Height: size.Height,
	})
}

