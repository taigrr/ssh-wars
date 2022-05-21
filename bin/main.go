package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taigrr/ssh-wars/asciimation"
)

func main() {
	var m asciimation.Model
	m.FrameSet = asciimation.ParseFrames()
	m.Progress = asciimation.ModelProg{Progress: progress.New(progress.WithSolidFill("17")),
		MaxWidth: 65,
		Padding:  2}
	m.Help = asciimation.NewHelpModel()
	m.Speed = 15

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

}
