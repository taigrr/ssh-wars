package asciimation

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestParseFrames(t *testing.T) {
	frames := parseFrames()
	if len(frames) == 0 {
		t.Fatal("parseFrames returned no frames")
	}
	for i, f := range frames {
		if len(f.lines) == 0 {
			t.Errorf("frame %d has no lines", i)
		}
		if len(f.lines) > viewportY {
			t.Errorf("frame %d has %d lines, expected at most %d", i, len(f.lines), viewportY)
		}
	}
}

func TestFrameRender(t *testing.T) {
	renderer := lipgloss.DefaultRenderer()
	yellow := renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffc500"))
	blue := renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("#174ea6"))

	f := Frame{
		lines:      []string{"Hello, World!"},
		frameCount: 1,
		index:      0,
	}
	rendered := f.Render(yellow, blue)
	if !strings.Contains(rendered, "Hello, World!") {
		t.Error("rendered frame does not contain expected text")
	}
	if !strings.Contains(rendered, "=") {
		t.Error("rendered frame missing border")
	}
}

func TestNewModel(t *testing.T) {
	renderer := lipgloss.DefaultRenderer()
	m := New(renderer)
	if m.currentFrame != 0 {
		t.Errorf("expected initial frame 0, got %d", m.currentFrame)
	}
	if m.paused {
		t.Error("expected model to start unpaused")
	}
}

func TestModelUpdate_Quit(t *testing.T) {
	renderer := lipgloss.DefaultRenderer()
	m := New(renderer)
	m.Help = NewHelpModel(renderer)

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if updated == nil {
		t.Fatal("Update returned nil model")
	}
	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
}

func TestModelUpdate_SpeedControls(t *testing.T) {
	renderer := lipgloss.DefaultRenderer()
	m := New(renderer)
	m.Help = NewHelpModel(renderer)
	m.Speed = 5

	// Speed up
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	um := updated.(Model)
	if um.Speed != 6 {
		t.Errorf("expected speed 6 after up, got %d", um.Speed)
	}

	// Speed down
	updated, _ = um.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	um = updated.(Model)
	if um.Speed != 5 {
		t.Errorf("expected speed 5 after down, got %d", um.Speed)
	}

	// Speed floor
	um.Speed = 1
	updated, _ = um.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	um = updated.(Model)
	if um.Speed != 1 {
		t.Errorf("expected speed to stay at 1, got %d", um.Speed)
	}
}

func TestModelUpdate_Navigation(t *testing.T) {
	renderer := lipgloss.DefaultRenderer()
	m := New(renderer)
	m.Help = NewHelpModel(renderer)

	// Forward
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	um := updated.(Model)
	if um.currentFrame != 1 {
		t.Errorf("expected frame 1 after right, got %d", um.currentFrame)
	}

	// Back
	updated, _ = um.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	um = updated.(Model)
	if um.currentFrame != 0 {
		t.Errorf("expected frame 0 after left, got %d", um.currentFrame)
	}

	// Can't go below 0
	updated, _ = um.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	um = updated.(Model)
	if um.currentFrame != 0 {
		t.Errorf("expected frame to stay at 0, got %d", um.currentFrame)
	}

	// Jump to end
	updated, _ = um.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	um = updated.(Model)
	if um.currentFrame != len(frameSet)-1 {
		t.Errorf("expected last frame, got %d", um.currentFrame)
	}
}

func TestModelUpdate_NumberJump(t *testing.T) {
	renderer := lipgloss.DefaultRenderer()
	m := New(renderer)
	m.Help = NewHelpModel(renderer)

	// Jump to 50%
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
	um := updated.(Model)
	expected := (len(frameSet) - 1) * 5 / 10
	if um.currentFrame != expected {
		t.Errorf("expected frame %d for 50%% jump, got %d", expected, um.currentFrame)
	}
}

func TestModelUpdate_Pause(t *testing.T) {
	renderer := lipgloss.DefaultRenderer()
	m := New(renderer)
	m.Help = NewHelpModel(renderer)
	m.Speed = 15

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	um := updated.(Model)
	if !um.paused {
		t.Error("expected model to be paused after space")
	}

	updated, _ = um.Update(tea.KeyMsg{Type: tea.KeySpace})
	um = updated.(Model)
	if um.paused {
		t.Error("expected model to be unpaused after second space")
	}
}

func TestModelUpdate_WindowTooSmall(t *testing.T) {
	renderer := lipgloss.DefaultRenderer()
	m := New(renderer)
	m.Help = NewHelpModel(renderer)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 50, Height: 10})
	um := updated.(Model)
	if !um.tooSmall {
		t.Error("expected tooSmall to be true for small window")
	}
	if !um.paused {
		t.Error("expected paused to be true for small window")
	}

	view := um.View()
	if !strings.Contains(view, "too small") {
		t.Error("expected 'too small' message in view")
	}
}

func TestModelView(t *testing.T) {
	renderer := lipgloss.DefaultRenderer()
	m := New(renderer)
	m.Help = NewHelpModel(renderer)

	view := m.View()
	if len(view) == 0 {
		t.Error("View returned empty string")
	}
}
