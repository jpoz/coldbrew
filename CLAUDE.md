# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`trmnl` is a Go library for building terminal applications with flexbox-like layout systems. It provides a component-based architecture for creating complex terminal UIs with proper layout, styling, and rendering.

## Commands

### Building and Running
- `go run examples/border/main.go` - Run the border layout example (demonstrates responsive rendering)
- `go build .` - Build the library
- `go test ./...` - Run tests (when they exist)
- `go mod tidy` - Clean up module dependencies

### Development
- `go run <example_path>` - Run specific examples to test functionality
- `gofmt -w .` - Format all Go code
- `go vet ./...` - Run Go static analysis

## Architecture

### Core Components
- **Component Interface** (`component.go`) - Base interface for all renderable elements with `Render()`, `RenderToBuffer()`, `GetMinSize()`, `GetStyle()`
- **Terminal** (`terminal.go`) - Main rendering engine with responsive rendering capabilities
  - `Render(component, size)` - Render with explicit size
  - `RenderResponsive(component)` - Auto-detect and use full terminal size
  - `RenderFullWidth(component, height)` - Use full terminal width with fixed height
  - `RenderFullHeight(component, width)` - Use full terminal height with fixed width
  - `GetSize()` - Detect current terminal dimensions
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

## File Organization
- Root package contains core library components
- `examples/` directory contains demonstration applications
- Each component type has its own file (text.go, flex.go, etc.)
- Shared types and utilities in dedicated files (types.go, color.go)