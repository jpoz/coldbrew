# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`trmnl` is a Go library for building terminal applications with flexbox-like layout systems. It provides a component-based architecture for creating complex terminal UIs with proper layout, styling, and rendering.

## Commands

### Building and Running
- `go run examples/border/main.go` - Run the border layout example (demonstrates responsive rendering)
- `go run examples/counter/main.go` - Run the counter example (demonstrates Elm architecture)
- `go run examples/clock/main.go` - Run the clock example (demonstrates subscriptions)
- `go run examples/async/main.go` - Run the async operations example (demonstrates commands)
- `go run examples/cursor/main.go` - Run the cursor control example (demonstrates cursor management)
- `go run examples/immediate/main.go` - Run the immediate input example (demonstrates raw mode)
- `go run examples/simple-immediate/main.go` - Run simple immediate input test
- `go build .` - Build the library
- `go test ./...` - Run tests (when they exist)
- `go mod tidy` - Clean up module dependencies

### Development
- `go run <example_path>` - Run specific examples to test functionality
- `gofmt -w .` - Format all Go code
- `go vet ./...` - Run Go static analysis

## Architecture

### Elm Architecture (The Elm Architecture pattern)
`trmnl` implements the Elm Architecture for interactive applications:
- **Model** - Application state and behavior interface with `Init()`, `Update()`, `View()`, `Subscriptions()`
- **Messages** - Events that update the model (keyboard input, timer ticks, async results)
- **Commands** - Side effects that produce messages asynchronously (`Cmd` type)
- **Subscriptions** - Continuous event streams (timers, external events)
- **Program** - Runtime that manages the Model-Update-View cycle

**Program Flow:**
1. `Init()` creates initial model and optional command
2. Input/subscriptions generate messages
3. `Update(msg)` processes messages and returns new model + optional command  
4. `View()` renders current model to component tree
5. Terminal renders component tree to console
6. Commands run asynchronously and send messages back

### Core Components
- **Component Interface** (`component.go`) - Base interface for all renderable elements with `Render()`, `RenderToBuffer()`, `GetMinSize()`, `GetStyle()`
- **Program** (`program.go`) - Elm architecture runtime
  - `NewProgram(model)` - Create program with initial model
  - `Run()` - Start the Model-Update-View loop 
  - `Send(msg)` - Send message to program
  - `WithCursorHidden(bool)` - Configure cursor visibility (default: hidden)
  - `WithRawMode(bool)` - Enable immediate character input without Enter (default: false)
  - Built-in commands: `Quit()`, `Delay()`, `Batch()`, `Tick()`
  - Built-in subscriptions: `Every(duration, msgFunc)`
- **Terminal** (`terminal.go`) - Main rendering engine with responsive rendering capabilities
  - `Render(component, size)` - Render with explicit size
  - `RenderResponsive(component)` - Auto-detect and use full terminal size
  - `RenderFullWidth(component, height)` - Use full terminal width with fixed height
  - `RenderFullHeight(component, width)` - Use full terminal height with fixed width
  - `GetSize()` - Detect current terminal dimensions
  - `HideCursor()` / `ShowCursor()` - Control cursor visibility
  - `MoveCursor(row, col)` / `MoveCursorHome()` - Cursor positioning
- **RenderBuffer** (`color.go`) - Efficient rendering system that separates content from color information
- **Text Component** (`text.go`) - Basic text rendering with styling support
- **FlexContainer** (`flex.go`) - Flexbox layout implementation with direction, justify, align properties

### Type System
- **Style Types** (`types.go`) - Colors, layout enums (Row/Column, JustifyStart/Center/SpaceBetween, etc.), Box model
- **Layout Properties** - FlexDirection, JustifyContent, AlignItems, flex-grow
- **Styling** - Borders (normal/rounded), padding, margins, ANSI colors

### Rendering Pipeline
1. Components calculate minimum required size via `GetMinSize()`
2. FlexContainer performs layout calculations using flexbox algorithm
3. Components render to RenderBuffer with separate content and color channels
4. Terminal outputs final buffer to console with ANSI escape sequences

### Key Design Patterns
- **Builder Pattern** - FlexContainer uses method chaining for configuration
- **Component Tree** - Hierarchical structure where containers hold child components
- **Separation of Concerns** - Content rendering separated from color/styling via RenderBuffer
- **Responsive Layout** - flex-grow allows components to expand and fill available space

### Input Modes
- **Line Mode (default)** - Traditional line-buffered input requiring Enter key
- **Raw Mode** - Immediate character input without Enter, supports:
  - Individual character detection
  - Arrow keys and special keys (space, tab, backspace, etc.)
  - Ctrl combinations (Ctrl+C, Ctrl+D)
  - Automatic terminal state restoration
  - Graceful fallback to line mode if raw mode unavailable

## File Organization
- Root package contains core library components
- `examples/` directory contains demonstration applications
- Each component type has its own file (text.go, flex.go, etc.)
- Shared types and utilities in dedicated files (types.go, color.go)