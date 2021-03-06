package asciimation

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Help    key.Binding
	Quit    key.Binding
	Numbers key.Binding
	Space   key.Binding
	G       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Space, k.Up, k.Down},
		{k.Left, k.Right, k.Help},
		{k.Numbers, k.G, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "increase speed"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "decrease speed"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "back one frame"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "forward one frame"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Numbers: key.NewBinding(
		key.WithHelp("0-9", "jump to position"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "play/pause"),
	),
	G: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "jump to end"),
	),
}

type HelpModel struct {
	keys             keyMap
	help             help.Model
	descriptionStyle lipgloss.Style
	keyStyle         lipgloss.Style
	lastKey          string
	quitting         bool
}

func NewHelpModel() HelpModel {
	h := HelpModel{
		keys:             keys,
		help:             help.New(),
		descriptionStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#DC3E58")),
		keyStyle:         lipgloss.NewStyle().Foreground(lipgloss.Color("#ffc500")),
	}
	h.help.Styles.ShortKey = h.keyStyle
	h.help.Styles.FullKey = h.keyStyle
	h.help.Styles.FullDesc = h.descriptionStyle
	h.help.Styles.ShortDesc = h.descriptionStyle
	return h
}

func (m HelpModel) UpdateDoeFoot(df lipgloss.DoeFoot) HelpModel {
	m.help = m.help.UpdateDoeFoot(df)
	return m
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Help) {
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

func (m HelpModel) View() string {
	out := m.help.View(m.keys)
	if !m.help.ShowAll {
		out += "\n\n"
	}
	return out
}
