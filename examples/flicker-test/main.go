package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jpoz/trmnl"
)

// FlickerTestModel demonstrates rapid updates to test flicker-free rendering
type FlickerTestModel struct {
	counter     int
	progress    int
	maxProgress int
	direction   int
}

func (m FlickerTestModel) Init() (trmnl.Model, trmnl.Cmd) {
	return m, trmnl.Tick(50*time.Millisecond, func(t time.Time) trmnl.Msg {
		return TickMsg{Time: t}
	})
}

type TickMsg struct {
	Time time.Time
}

func (m FlickerTestModel) Update(msg trmnl.Msg) (trmnl.Model, trmnl.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		// Update counter
		m.counter++
		
		// Update bouncing progress bar
		m.progress += m.direction
		if m.progress >= m.maxProgress {
			m.direction = -1
		} else if m.progress <= 0 {
			m.direction = 1
		}
		
		// Continue ticking
		return m, trmnl.Tick(50*time.Millisecond, func(t time.Time) trmnl.Msg {
			return TickMsg{Time: t}
		})
	case trmnl.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, trmnl.Quit()
		}
	}
	return m, nil
}

func (m FlickerTestModel) View() trmnl.Component {
	// Create a rapidly changing display to test for flicker
	title := trmnl.NewText("🔥 Flicker Test - Rapid Updates")
	title.Style.BorderColor = trmnl.ColorYellow
	title.Style.Border = true
	title.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	counterText := trmnl.NewText(fmt.Sprintf("Counter: %d", m.counter))
	counterText.Style.BorderColor = trmnl.ColorCyan
	counterText.Style.Border = true
	counterText.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Create a bouncing progress bar
	progressBar := m.createProgressBar()
	
	// Create rapidly changing content
	timeInfo := trmnl.NewText(fmt.Sprintf("Time: %s", time.Now().Format("15:04:05.000")))
	timeInfo.Style.BorderColor = trmnl.ColorGreen
	timeInfo.Style.Border = true
	timeInfo.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Bouncing dots animation
	dots := strings.Repeat(".", (m.counter%10)+1)
	animText := trmnl.NewText(fmt.Sprintf("Animation%s", dots))
	animText.Style.BorderColor = trmnl.ColorMagenta
	animText.Style.Border = true
	animText.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	instructions := trmnl.NewText("Press 'q' or Ctrl+C to exit. Watch for flicker - there should be none!")
	instructions.Style.BorderColor = trmnl.ColorWhite
	instructions.Style.Border = true
	instructions.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	// Layout everything in a column
	container := trmnl.NewFlexContainer().
		SetDirection(trmnl.Column).
		SetJustify(trmnl.JustifyStart).
		SetAlign(trmnl.AlignCenter).
		SetPadding(trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}).
		AddChild(title).
		AddChild(counterText).
		AddChild(progressBar).
		AddChild(timeInfo).
		AddChild(animText).
		AddChild(instructions)
	
	return container
}

func (m FlickerTestModel) createProgressBar() trmnl.Component {
	width := 40
	filled := (m.progress * width) / m.maxProgress
	
	var bar strings.Builder
	bar.WriteString("[")
	
	for i := 0; i < width; i++ {
		if i <= filled {
			bar.WriteString("█")
		} else {
			bar.WriteString("░")
		}
	}
	
	bar.WriteString("]")
	bar.WriteString(fmt.Sprintf(" %d%%", (m.progress*100)/m.maxProgress))
	
	progressText := trmnl.NewText(bar.String())
	progressText.Style.BorderColor = trmnl.ColorBlue
	progressText.Style.Border = true
	progressText.Style.Padding = trmnl.Box{Top: 1, Right: 2, Bottom: 1, Left: 2}
	
	return progressText
}

func (m FlickerTestModel) Subscriptions() []trmnl.Sub {
	return nil
}

func main() {
	model := FlickerTestModel{
		counter:     0,
		progress:    0,
		maxProgress: 20,
		direction:   1,
	}
	
	program := trmnl.NewProgram(model).
		WithCursorHidden(true).
		WithRawMode(true)
	
	if err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}