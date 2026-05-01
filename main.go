package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/spf13/cobra"
	"github.com/taigrr/ssh-wars/asciimation"
)

func main() {
	var host string
	var port int

	command := &cobra.Command{
		Use:   "ssh-wars",
		Short: "Serve the Star Wars ASCII animation over SSH",
		RunE: func(cmd *cobra.Command, args []string) error {
			server, err := wish.NewServer(
				wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
				wish.WithHostKeyPath(".ssh/term_info_ed25519"),
				wish.WithMiddleware(
					bm.Middleware(teaHandler),
					activeterm.Middleware(),
					lm.Middleware(),
				),
			)
			if err != nil {
				return err
			}

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer signal.Stop(done)

			log.Printf("Starting SSH server on %s:%d", host, port)
			go func() {
				if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, ssh.ErrServerClosed) {
					log.Printf("server error: %v", serveErr)
					done <- syscall.SIGTERM
				}
			}()

			<-done
			log.Println("Stopping SSH server")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if shutdownErr := server.Shutdown(ctx); shutdownErr != nil && !errors.Is(shutdownErr, ssh.ErrServerClosed) {
				return shutdownErr
			}
			return nil
		},
	}

	command.Flags().StringVar(&host, "host", "0.0.0.0", "Host to listen on")
	command.Flags().IntVar(&port, "port", 2222, "Port to listen on")

	if err := fang.Execute(context.Background(), command); err != nil {
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
	model := asciimation.NewDefaultModel(renderer)

	return model, []tea.ProgramOption{tea.WithAltScreen()}
}
