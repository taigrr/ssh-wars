package asciimation

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestModelProg_Init(t *testing.T) {
	m := ModelProg{}
	if cmd := m.Init(); cmd != nil {
		t.Error("expected Init to return nil")
	}
}

func TestModelProg_Update_WindowResize(t *testing.T) {
	m := ModelProg{
		Progress: progress.New(progress.WithColors(lipgloss.Color("#174ea6"))),
		MaxWidth: 65,
		Padding:  2,
	}

	// Small window: width should be calculated as Width - Padding*2 - 4
	um, _ := m.Update(tea.WindowSizeMsg{Width: 50, Height: 40})
	expected := 50 - 2*2 - 4
	if um.Progress.Width() != expected {
		t.Errorf("expected width %d, got %d", expected, um.Progress.Width())
	}

	// Very small window: width should never go negative.
	um, _ = m.Update(tea.WindowSizeMsg{Width: 6, Height: 40})
	if um.Progress.Width() != 0 {
		t.Errorf("expected width clamped at 0, got %d", um.Progress.Width())
	}

	// Large window: width should be capped at MaxWidth
	um, _ = m.Update(tea.WindowSizeMsg{Width: 200, Height: 40})
	if um.Progress.Width() != 65 {
		t.Errorf("expected width capped at 65, got %d", um.Progress.Width())
	}
}

func TestModelProg_View(t *testing.T) {
	m := ModelProg{
		Progress: progress.New(progress.WithColors(lipgloss.Color("#174ea6"))),
		MaxWidth: 65,
		Padding:  2,
		percent:  0.5,
	}
	view := m.View()
	if !strings.Contains(view, "\n") {
		t.Error("expected view to contain newlines")
	}
}
