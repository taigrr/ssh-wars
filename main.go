package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//go:embed starwars.ascii
var asciiString string

var frames []Frame

const viewportY = 13

type Frame struct {
	lines      []string
	frameCount int
}

func (f Frame) String() string {
	return strings.Join(f.lines, "\n")
}

func init() {
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
}

func main() {
	for _, f := range frames {
		fmt.Println(f)
		time.Sleep(time.Second * time.Duration(f.frameCount) / 15)
	}
}
