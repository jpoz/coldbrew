package brew

// Color represents terminal colors using ANSI escape codes
type Color string

const (
	ColorDefault Color = ""
	ColorBlack   Color = "\033[30m"
	ColorRed     Color = "\033[31m"
	ColorGreen   Color = "\033[32m"
	ColorYellow  Color = "\033[33m"
	ColorBlue    Color = "\033[34m"
	ColorMagenta Color = "\033[35m"
	ColorCyan    Color = "\033[36m"
	ColorWhite   Color = "\033[37m"

	// Bright colors
	ColorBrightBlack   Color = "\033[90m"
	ColorBrightRed     Color = "\033[91m"
	ColorBrightGreen   Color = "\033[92m"
	ColorBrightYellow  Color = "\033[93m"
	ColorBrightBlue    Color = "\033[94m"
	ColorBrightMagenta Color = "\033[95m"
	ColorBrightCyan    Color = "\033[96m"
	ColorBrightWhite   Color = "\033[97m"

	// Color reset
	ColorReset Color = "\033[0m"
)

// Direction represents flex direction
type Direction int

const (
	Row Direction = iota
	Column
)

// Justify represents justify-content alignment
type Justify int

const (
	JustifyStart Justify = iota
	JustifyCenter
	JustifyEnd
	JustifySpaceBetween
	JustifySpaceAround
)

// Align represents align-items alignment
type Align int

const (
	AlignStart Align = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// Box represents spacing around content
type Box struct {
	Top, Right, Bottom, Left int
}

// Size represents dimensions
type Size struct {
	Width, Height int
}

// Position represents coordinates
type Position struct {
	X, Y int
}

// Style contains visual styling properties
type Style struct {
	Border        bool
	BorderChar    rune
	BorderColor   Color
	RoundedBorder bool
	Padding       Box
	Margin        Box
	BgChar        rune
}

