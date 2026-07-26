package asciimation

import (
	_ "embed"
	"strconv"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

//go:embed intro.ascii
var introString string

//go:embed starwars.ascii
var asciiString string

var onceFrames sync.Once

var frameSet []Frame

const viewportY = 13
const longAgoFrame = 49
const scrawlStart = 51
const scrawlEnd = 111

// Frame geometry: each rendered line is "||" + 67 content columns + "||".
const (
	frameContentWidth = 67
	frameBorderWidth  = 71
	frameLineWidth    = 72 // border width plus trailing newline

	minWindowWidth  = 72
	minWindowHeight = 22
)

type TickMsg struct{}

type Model struct {
	Progress     ModelProg
	Help         HelpModel
	Speed        int
	AltScreen    bool
	yellow       lipgloss.Style
	blue         lipgloss.Style
	currentFrame int
	paused       bool
	tooSmall     bool
}

type Frame struct {
	lines      []string
	frameCount int
	index      int
}

func New() Model {
	m := Model{
		yellow: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffc500")),
		blue:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#174ea6")),
	}
	onceFrames.Do(func() {
		frameSet = parseFrames()
	})
	return m
}

func (f Frame) Render(yellow, blue lipgloss.Style) string {
	localBorder := yellow.Render(strings.Repeat("=", frameBorderWidth))
	edge := yellow.Render("||")

	var sb strings.Builder
	sb.Grow((len(f.lines) + 2) * frameLineWidth)
	sb.WriteString(localBorder)
	sb.WriteString("\n")
	for _, l := range f.lines {
		sb.WriteString(edge)
		length := len(l)
		if f.index == longAgoFrame {
			l = blue.Render(l)
		} else if f.index < scrawlEnd && f.index >= scrawlStart {
			l = yellow.Render(l)
		}
		sb.WriteString(l)
		sb.WriteString(strings.Repeat(" ", max(0, frameContentWidth-length)))
		sb.WriteString(edge)
		sb.WriteString("\n")
	}
	sb.WriteString(localBorder)
	return sb.String()
}

func parseFrames() []Frame {
	var frames []Frame
	f := Frame{index: 0}
	lines := strings.Split(introString, "\n")
	asciiString = strings.ReplaceAll(asciiString, "\\'", "'")
	asciiString = strings.ReplaceAll(asciiString, "\"", "\\\"")
	lines = append(lines, strings.Split(asciiString, "\\n")...)
	for i, l := range lines {
		if i%(viewportY+1) == 0 {
			f = Frame{index: i / (viewportY + 1)}
			f.frameCount, _ = strconv.Atoi(l)
			continue
		}
		u, err := strconv.Unquote("\"" + l + "\"")
		if err != nil {
			u = l
		}
		f.lines = append(f.lines, u)
		if i%(viewportY+1) == viewportY {
			frames = append(frames, f)
		}
	}
	return frames
}

func (m Model) View() tea.View {
	var content string
	if m.tooSmall {
		content = "Window is too small for player.\nPlease try resizing your window."
	} else {
		content = frameSet[m.currentFrame].Render(m.yellow, m.blue) + "\n" + m.Progress.View() + m.Help.View() + "\n"
	}
	v := tea.NewView(content)
	v.AltScreen = m.AltScreen
	return v
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second*time.Duration(frameSet[m.currentFrame].frameCount)/time.Duration(m.Speed), func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TickMsg:
		if m.paused {
			break
		}
		if m.currentFrame < len(frameSet)-1 {
			m.currentFrame++
			cmd = m.tick()
		} else {
			m.paused = true
		}
	case tea.WindowSizeMsg:
		if msg.Width < minWindowWidth || msg.Height < minWindowHeight {
			m.paused = true
			m.tooSmall = true
		} else if m.tooSmall {
			m.paused = false
			m.tooSmall = false
			cmd = m.tick()
		}
		m.Help, _ = m.Help.Update(msg)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "right", "l":
			if m.currentFrame < len(frameSet)-1 {
				m.currentFrame++
			}
		case "up", "k":
			m.Speed++
		case "down", "j":
			if m.Speed > 1 {
				m.Speed--
			}
		case "left", "h":
			if m.currentFrame > 0 {
				m.currentFrame--
			}
		case "G":
			m.currentFrame = len(frameSet) - 1
		case "g":
			m.currentFrame = 0
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			num, _ := strconv.Atoi(msg.String())
			m.currentFrame = (len(frameSet) - 1) * num / 10
		case "space":
			m.paused = !m.paused
			if !m.paused {
				return m, m.tick()
			}
			return m, nil
		default:
			m.Help, _ = m.Help.Update(msg)
		}
	}
	m.Progress.percent = progressPercent(m.currentFrame, len(frameSet))
	m.Progress, _ = m.Progress.Update(msg)

	return m, cmd
}

func progressPercent(currentFrame, frameCount int) float64 {
	if frameCount <= 1 {
		return 1
	}

	percent := float64(currentFrame) / float64(frameCount-1)
	if percent < 0 {
		return 0
	}
	if percent > 1 {
		return 1
	}
	return percent
}
