package trmnl

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Type aliases for bubbletea compatibility
type Msg = tea.Msg
type Cmd = tea.Cmd
type Model = tea.Model

// Sub represents a subscription that can send messages continuously
type Sub func(ctx context.Context, send func(Msg))

// TrmnlModel extends tea.Model with subscriptions for advanced functionality
type TrmnlModel interface {
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
		hideCursor: true,  // Default to hiding cursor for TUI apps
		rawMode:    false, // Default to line-buffered input
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

// Run starts the program and blocks until it exits
func (p *Program) Run() error {
	// Setup cursor visibility
	if p.hideCursor {
		p.terminal.HideCursor()
		// Ensure cursor is restored on exit
		defer p.terminal.ShowCursor()
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

	// Main message loop
	for {
		select {
		case msg := <-p.msgChan:
			// Handle quit message
			if _, isQuit := msg.(QuitMsg); isQuit {
				p.cancel()
				return nil
			}

			// Handle batch messages
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
			return nil
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

// startSubscriptions starts all subscriptions from the model if it supports them
func (p *Program) startSubscriptions() {
	// Check if model supports subscriptions
	if trmnlModel, ok := p.model.(TrmnlModel); ok {
		p.subs = trmnlModel.Subscriptions()

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

