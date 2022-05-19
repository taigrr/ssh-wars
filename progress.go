package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type modelProg struct {
	maxWidth int
	padding  int
	percent  float64
	progress progress.Model
}

func (_ modelProg) Init() tea.Cmd {
	return nil
}

func (m modelProg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - m.padding*2 - 4
		if m.progress.Width > m.maxWidth {
			m.progress.Width = m.maxWidth
		}
		return m, nil

	default:
		return m, nil
	}
}

func (e modelProg) View() string {
	pad := strings.Repeat(" ", e.padding)
	return "\n" +
		pad + e.progress.ViewAs(e.percent) + "\n\n"
}
