package asciimation

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelProg struct {
	MaxWidth int
	Padding  int
	percent  float64
	Progress progress.Model
}

func (_ ModelProg) Init() tea.Cmd {
	return nil
}

func (m ModelProg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Progress.Width = msg.Width - m.Padding*2 - 4
		if m.Progress.Width > m.MaxWidth {
			m.Progress.Width = m.MaxWidth
		}
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
