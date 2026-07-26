package asciimation

import (
	"strings"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
)

type ModelProg struct {
	MaxWidth int
	Padding  int
	percent  float64
	Progress progress.Model
}

func (m ModelProg) Update(msg tea.Msg) (ModelProg, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w := max(0, msg.Width-m.Padding*2-4)
		w = min(w, m.MaxWidth)
		m.Progress.SetWidth(w)
		return m, nil
	default:
		return m, nil
	}
}

func (e ModelProg) View() string {
	pad := strings.Repeat(" ", e.Padding)
	return "\n" +
		pad + e.Progress.ViewAs(e.percent) + "\n\n"
}
