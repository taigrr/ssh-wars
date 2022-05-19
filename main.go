package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed starwars.ascii
var asciiString string

const viewportY = 13

type TickMsg struct {
	Time time.Time
	tag  int
	ID   int
}

type model struct {
	frameSet     []Frame
	speed        int
	currentFrame int
	paused       bool
}

type Frame struct {
	lines      []string
	frameCount int
}

func (f Frame) String() string {
	return strings.Join(f.lines, "\n")
}

func parseFrames() []Frame {
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

func (m model) View() string {
	return m.frameSet[m.currentFrame].String()
}

func (m model) Init() tea.Cmd {
	return m.tick(m.currentFrame, m.currentFrame)
}

func (m model) tick(id, tag int) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(m.frameSet[m.currentFrame].frameCount)/15, func(t time.Time) tea.Msg {
		return TickMsg{
			Time: t,
			ID:   id,
			tag:  tag,
		}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if m.paused {
			break
		}
		if m.currentFrame < len(m.frameSet)-1 {
			m.currentFrame++
			return m, m.tick(m.currentFrame, m.currentFrame)
		} else {
			m.paused = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l":
			if m.currentFrame < len(m.frameSet)-1 {
				m.currentFrame++
			}
		case "down", "j":
			m.speed--
		case "left", "h":
			if m.currentFrame > 0 {
				m.currentFrame--
			}
		case "G":
			m.currentFrame = len(m.frameSet) - 1
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			num, _ := strconv.Atoi(msg.String())
			m.currentFrame = len(m.frameSet) - 1
			m.currentFrame = m.currentFrame * num / 10
		case " ":
			m.paused = !m.paused
			return m, m.tick(m.currentFrame, m.currentFrame)
		}
	}
	return m, nil
}

func main() {
	var m model
	m.frameSet = parseFrames()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

}
