package asciimation

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestNewHelpModel(t *testing.T) {
	h := NewHelpModel()
	if h.help.ShowAll {
		t.Error("expected help to start collapsed")
	}
}

func TestHelpModel_ToggleHelp(t *testing.T) {
	h := NewHelpModel()

	um, _ := h.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	if !um.help.ShowAll {
		t.Error("expected help to be expanded after ?")
	}

	um, _ = um.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	if um.help.ShowAll {
		t.Error("expected help to be collapsed after second ?")
	}
}

func TestHelpModel_View(t *testing.T) {
	h := NewHelpModel()
	view := h.View()
	if len(view) == 0 {
		t.Error("HelpModel View returned empty string")
	}
}

func TestHelpModel_WindowResize(t *testing.T) {
	h := NewHelpModel()

	um, _ := h.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if um.help.Width() != 120 {
		t.Errorf("expected help width 120, got %d", um.help.Width())
	}
}

func TestKeyMap_ShortHelp(t *testing.T) {
	bindings := keys.ShortHelp()
	if len(bindings) != 2 {
		t.Errorf("expected 2 short help bindings, got %d", len(bindings))
	}
}

func TestKeyMap_FullHelp(t *testing.T) {
	groups := keys.FullHelp()
	if len(groups) != 3 {
		t.Errorf("expected 3 full help groups, got %d", len(groups))
	}
}
