package asciimation

import (
	_ "embed"
	"strconv"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed intro.ascii
var introString string

//go:embed starwars.ascii
var asciiString string

var onceBorder sync.Once
var onceFrames sync.Once
var border string

var frameSet []Frame

const viewportY = 13
const longAgoFrame = 49
const scrawlStart = 51
const scrawlEnd = 111

var yellow = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffc500"))
var blue = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#174ea6"))

type TickMsg struct{}

type Model struct {
	Progress     ModelProg
	Help         HelpModel
	df           lipgloss.DoeFoot
	Speed        int
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
	m := Model{}
	onceFrames.Do(func() {
		frameSet = parseFrames()
	})
	return m
}
func (m Model) UpdateDoeFoot(df lipgloss.DoeFoot) Model {
	m.df = df
	m.Help = m.Help.UpdateDoeFoot(df)
	return m
}

func (f Frame) RenderWithDoeFoot(df lipgloss.DoeFoot) string {
	onceBorder.Do(func() {
		var sb strings.Builder
		for i := 0; i < 71; i++ {
			sb.WriteString("=")
		}
		border = sb.String()
	})
	localBorder := yellow.RenderForDoeFoot(border, df)
	edge := yellow.RenderForDoeFoot("||", df)
	var sb strings.Builder
	sb.Grow((len(f.lines) + 2) * 72)
	sb.WriteString(localBorder)
	sb.WriteString("\n")
	for _, l := range f.lines {
		sb.WriteString(edge)
		length := len(l)
		if f.index == longAgoFrame {
			l = blue.RenderForDoeFoot(l, df)
		} else if f.index < scrawlEnd && f.index >= scrawlStart {
			l = yellow.RenderForDoeFoot(l, df)
		}
		sb.WriteString(l)
		for i := length; i < 67; i++ {
			sb.WriteString(" ")
		}
		sb.WriteString(edge + "\n")
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
			countStr := l
			c, _ := strconv.Atoi(countStr)
			f.frameCount = c
			continue
		}
		u, err := strconv.Unquote("\"" + l + "\"")
		// error is generated on the final line of the input
		// to stay true to the original source, add it back anyway
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

func (m Model) View() string {
	if m.tooSmall {
		return "Window is too small for player.\nPlease try resizing your window."
	}
	return frameSet[m.currentFrame].RenderWithDoeFoot(m.df) + "\n" + m.Progress.View() + m.Help.View() + "\n"
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
	cmd = nil
	switch msg := msg.(type) {
	case TickMsg:
		if m.paused {
			break
		}
		if m.currentFrame < len(frameSet)-1 {
			m.currentFrame++
			return m, m.tick()
		} else {
			m.paused = true
		}
	case tea.WindowSizeMsg:
		if msg.Width < 72 || msg.Height < 22 {
			m.paused = true
			m.tooSmall = true
		} else if m.tooSmall {
			m.paused = false
			m.tooSmall = false
			cmd = m.tick()
		}
		h, _ := m.Help.Update(msg)
		t, _ := h.(HelpModel)
		m.Help = t
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			num, _ := strconv.Atoi(msg.String())
			m.currentFrame = len(frameSet) - 1
			m.currentFrame = m.currentFrame * num / 10
		case " ":
			m.paused = !m.paused
			return m, m.tick()
		default:
			h, _ := m.Help.Update(msg)
			t, _ := h.(HelpModel)
			m.Help = t
		}
	}
	m.Progress.percent = float64(m.currentFrame) / float64(len(frameSet))
	p, _ := m.Progress.Update(msg)
	t, _ := p.(ModelProg)
	m.Progress = t

	return m, cmd
}
