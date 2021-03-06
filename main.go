package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/muesli/termenv"
	"github.com/taigrr/ssh-wars/asciimation"
)

var host = flag.String("host", "0.0.0.0", "host to listen on")
var port = flag.Int("port", 2222, "port to listen on")

func main() {
	flag.Parse()
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", *host, *port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", *host, *port)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		fmt.Println("no active terminal, skipping")
		return nil, nil
	}
	m := asciimation.New()
	m.Progress = asciimation.ModelProg{Progress: progress.New(progress.WithSolidFill("#174ea6"),
		progress.WithColorProfile(getSSHTermInfo(pty))),
		MaxWidth: 65,
		Padding:  2}
	m.Help = asciimation.NewHelpModel()
	df := lipgloss.NewDoeFoot()
	df = df.SetColorProfile(getSSHTermInfo(pty))
	m = m.UpdateDoeFoot(df)
	m.Speed = 15

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func getSSHTermInfo(pty ssh.Pty) termenv.Profile {
	termOut := termenv.NewOutputWithProfileEnv(nil, pty.Term, "")
	return termOut.Profile
}
