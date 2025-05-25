# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`trmnl` is a Go library for building terminal applications with flexbox-like layout systems. It provides **full compatibility with Bubble Tea Models** while adding enhanced layout capabilities through string builder utilities. Any Bubble Tea Model can run directly in trmnl's Program without modification.

## Commands

### Building and Running
- `go run examples/simple/main.go` - Run basic bubbletea countdown example
- `go run examples/bubbletea-demo/main.go` - Run bubbletea counter with trmnl layouts (spinner demo)
- `go run examples/trmnl-subscriptions/main.go` - Run example showing trmnl's subscription system
- `go run examples/long-essay/main.go` - Run long scrollable essay with dynamic content updates
- `go run examples/key-test/main.go` - Test immediate key response with raw mode
- `go run examples/quit-test/main.go` - Test tea.Quit command functionality
- `go run examples/simple-quit/main.go` - Simple tea.Quit test
- `go build .` - Build the library
- `go test ./...` - Run tests (when they exist)
- `go mod tidy` - Clean up module dependencies

### Development
- `go run <example_path>` - Run specific examples to test functionality
- `gofmt -w .` - Format all Go code
- `go vet ./...` - Run Go static analysis

## Architecture

### Bubble Tea Compatibility
`trmnl` is **fully compatible with Bubble Tea Models**:
- **Model Interface** - Uses `tea.Model` interface directly (`Init()`, `Update()`, `View()`)
- **Type Aliases** - `trmnl.Msg`, `trmnl.Cmd`, `trmnl.Model` are aliases for `tea.Msg`, `tea.Cmd`, `tea.Model`
- **Enhanced Features** - Optional `TrmnlModel` interface adds `Subscriptions()` for advanced functionality
- **String-based Rendering** - `View() string` works directly with trmnl's string-based rendering pipeline

**Program Flow:**
1. `Init()` returns initial command (bubbletea standard)
2. Input generates bubbletea messages
3. `Update(msg)` processes messages and returns new model + optional command (bubbletea standard)
4. `View()` renders current model to string (bubbletea standard)
5. Terminal renders string to console with trmnl's enhanced capabilities
6. Commands run asynchronously (bubbletea standard)
7. Optional subscriptions provide continuous event streams (trmnl extension)

### Core Components
- **Program** (`program.go`) - Bubbletea-compatible runtime
  - `NewProgram(tea.Model)` - Create program with any bubbletea model
  - `Run()` - Start the Model-Update-View loop 
  - `Send(msg)` - Send message to program
  - `WithCursorHidden(bool)` - Configure cursor visibility (default: hidden)
  - `WithRawMode(bool)` - Enable immediate character input without Enter (default: false)
  - **Full bubbletea command support**: `tea.Quit`, `tea.Batch`, etc. work exactly as expected
  - **Quit handling**: Responds to both `tea.Quit` commands and `tea.QuitMsg` messages automatically
  - Optional trmnl subscriptions: `Every(duration, msgFunc)` via `TrmnlModel` interface
- **Terminal** (`terminal.go`) - Main rendering engine with responsive rendering capabilities
  - `RenderString(content)` - Render string content directly (primary method for bubbletea compatibility)
  - `Render(component, size)` - Render component with explicit size (legacy)
  - `RenderResponsive(component)` - Auto-detect and use full terminal size (legacy)
  - `GetSize()` - Detect current terminal dimensions
  - `HideCursor()` / `ShowCursor()` - Control cursor visibility
  - `MoveCursor(row, col)` / `MoveCursorHome()` - Cursor positioning
- **String Builders** (`builders.go`) - Utilities for creating layouts with strings (primary for bubbletea compatibility)
  - `BuildBorder(content, title)` - Create bordered text blocks
  - `BuildRoundedBorder(content, title)` - Create rounded bordered text blocks
  - `BuildFlex(direction, items...)` - Create flexible layouts combining multiple content blocks
  - `BuildCenter(content, width)` - Center content within specified width
  - `BuildPadding(content, padding)` - Add padding around content
  - `BuildProgress(current, total, width, title)` - Create progress bar displays
- **Component System** (legacy, maintained for internal use)
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
  - **Enable with**: `program.WithRawMode(true)` - Required for interactive examples

## File Organization
- Root package contains core library components
- `examples/` directory contains demonstration applications
- Each component type has its own file (text.go, flex.go, etc.)
- Shared types and utilities in dedicated files (types.go, color.go)