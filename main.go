package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
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
			activeterm.Middleware(),
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
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Fatalln(err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	_, _, active := s.Pty()
	if !active {
		fmt.Println("no active terminal, skipping")
		return nil, nil
	}

	renderer := bm.MakeRenderer(s)

	m := asciimation.New(renderer)
	m.Progress = asciimation.ModelProg{
		Progress: progress.New(progress.WithSolidFill("#174ea6")),
		MaxWidth: 65,
		Padding:  2,
	}
	m.Help = asciimation.NewHelpModel(renderer)
	m.Speed = 15

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
