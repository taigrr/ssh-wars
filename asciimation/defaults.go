package asciimation

import (
	"charm.land/bubbles/v2/progress"
	"charm.land/lipgloss/v2"
)

// NewDefaultModel constructs the animation model with the shared styling used by
// both the SSH server and the standalone client.
func NewDefaultModel() Model {
	model := New()
	model.Progress = ModelProg{
		Progress: progress.New(progress.WithColors(lipgloss.Color("#174ea6"))),
		MaxWidth: 65,
		Padding:  2,
	}
	model.Help = NewHelpModel()
	model.Speed = 15
	return model
}
