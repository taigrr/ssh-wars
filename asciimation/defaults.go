package asciimation

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

// NewDefaultModel constructs the animation model with the shared styling used by
// both the SSH server and the standalone client.
func NewDefaultModel(renderer *lipgloss.Renderer) Model {
	model := New(renderer)
	model.Progress = ModelProg{
		Progress: progress.New(progress.WithSolidFill("#174ea6")),
		MaxWidth: 65,
		Padding:  2,
	}
	model.Help = NewHelpModel(renderer)
	model.Speed = 15
	return model
}
