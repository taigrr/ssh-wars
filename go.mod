module github.com/taigrr/ssh-wars

go 1.18

replace github.com/charmbracelet/bubbles => ../bubbles

replace github.com/muesli/termenv => github.com/taigrr/termenv v0.11.2

replace github.com/charmbracelet/lipgloss => ../lipgloss

require (
	github.com/charmbracelet/bubbles v0.10.3
	github.com/charmbracelet/bubbletea v0.20.0
	github.com/charmbracelet/lipgloss v0.5.0
	github.com/charmbracelet/wish v0.4.0
	github.com/gliderlabs/ssh v0.3.4
	github.com/muesli/termenv v0.11.1-0.20220212125758-44cd13922739
)

require (
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/caarlos0/sshmarshal v0.1.0 // indirect
	github.com/charmbracelet/harmonica v0.2.0 // indirect
	github.com/charmbracelet/keygen v0.3.0 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/muesli/ansi v0.0.0-20211031195517-c9f0611b6c70 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
)
