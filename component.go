package trmnl

// Component interface that all renderable elements must implement
type Component interface {
	Render(size Size) []string
	RenderToBuffer(size Size) *RenderBuffer
	GetMinSize() Size
	GetStyle() Style
}