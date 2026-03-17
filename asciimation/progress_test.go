package asciimation

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModelProg_Init(t *testing.T) {
	m := ModelProg{}
	if cmd := m.Init(); cmd != nil {
		t.Error("expected Init to return nil")
	}
}

func TestModelProg_Update_WindowResize(t *testing.T) {
	m := ModelProg{
		Progress: progress.New(progress.WithSolidFill("#174ea6")),
		MaxWidth: 65,
		Padding:  2,
	}

	// Small window: width should be calculated as Width - Padding*2 - 4
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 50, Height: 40})
	um := updated.(ModelProg)
	expected := 50 - 2*2 - 4
	if um.Progress.Width != expected {
		t.Errorf("expected width %d, got %d", expected, um.Progress.Width)
	}

	// Large window: width should be capped at MaxWidth
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 200, Height: 40})
	um = updated.(ModelProg)
	if um.Progress.Width != 65 {
		t.Errorf("expected width capped at 65, got %d", um.Progress.Width)
	}
}

func TestModelProg_View(t *testing.T) {
	m := ModelProg{
		Progress: progress.New(progress.WithSolidFill("#174ea6")),
		MaxWidth: 65,
		Padding:  2,
		percent:  0.5,
	}
	view := m.View()
	if !strings.Contains(view, "\n") {
		t.Error("expected view to contain newlines")
	}
}
