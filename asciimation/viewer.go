package asciimation

import (
	_ "embed"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed starwars.ascii
var asciiString string

const viewportY = 13

type TickMsg struct {
}

type Model struct {
	Progress     ModelProg
	Help         HelpModel
	FrameSet     []Frame
	Speed        int
	currentFrame int
	paused       bool
	tooSmall     bool
}

type Frame struct {
	lines      []string
	frameCount int
}

func (f Frame) String() string {
	return strings.Join(f.lines, "\n")
}

func ParseFrames() []Frame {
	var frames []Frame
	asciiString = strings.ReplaceAll(asciiString, "\\'", "'")
	asciiString = strings.ReplaceAll(asciiString, "\\\\", "\\")
	lines := strings.Split(asciiString, "\\n")
	var f Frame
	for i, l := range lines {
		if i%(viewportY+1) == 0 {
			f = Frame{}
			countStr := l
			c, _ := strconv.Atoi(countStr)
			f.frameCount = c
			continue
		}
		f.lines = append(f.lines, l)
		if i%(viewportY+1) == viewportY {
			frames = append(frames, f)
		}
	}
	return frames
}

func (m Model) View() string {
	return m.FrameSet[m.currentFrame].String() + "\n" + m.Progress.View() + m.Help.View() + "\n"
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second*time.Duration(m.FrameSet[m.currentFrame].frameCount)/time.Duration(m.Speed), func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if m.paused {
			break
		}
		if m.currentFrame < len(m.FrameSet)-1 {
			m.currentFrame++
			return m, m.tick()
		} else {
			m.paused = true
		}
	case tea.WindowSizeMsg:
		if msg.Width < 68 {
			m.paused = true
			m.tooSmall = true
		} else if m.tooSmall {
			m.paused = false
		}
		h, _ := m.Help.Update(msg)
		t, _ := h.(HelpModel)
		m.Help = t
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l":
			if m.currentFrame < len(m.FrameSet)-1 {
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
			m.currentFrame = len(m.FrameSet) - 1
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			num, _ := strconv.Atoi(msg.String())
			m.currentFrame = len(m.FrameSet) - 1
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
	m.Progress.percent = float64(m.currentFrame) / float64(len(m.FrameSet))
	p, _ := m.Progress.Update(msg)
	t, _ := p.(ModelProg)
	m.Progress = t

	return m, nil
}
