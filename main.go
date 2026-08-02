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

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	bm "charm.land/wish/v2/bubbletea"
	lm "charm.land/wish/v2/logging"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/ssh"
	"github.com/spf13/cobra"
	"github.com/taigrr/ssh-wars/asciimation"
)

const (
	defaultHost        = "0.0.0.0"
	defaultPort        = 2222
	defaultHostKeyPath = ".ssh/term_info_ed25519"
	shutdownTimeout    = 30 * time.Second
)

type serverConfig struct {
	host        string
	port        int
	hostKeyPath string
}

func main() {
	if err := fang.Execute(context.Background(), newRootCommand()); err != nil {
		log.Fatalln(err)
	}
}

func newRootCommand() *cobra.Command {
	config := serverConfig{
		host:        defaultHost,
		port:        defaultPort,
		hostKeyPath: defaultHostKeyPath,
	}
	command := &cobra.Command{
		Use:   "ssh-wars",
		Short: "Serve the Star Wars ASCII animation over SSH",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(config)
		},
	}

	command.Flags().StringVar(&config.host, "host", defaultHost, "Host to listen on")
	command.Flags().IntVar(&config.port, "port", defaultPort, "Port to listen on")
	command.Flags().StringVar(&config.hostKeyPath, "host-key-path", defaultHostKeyPath, "SSH host key path")

	return command
}

func runServer(config serverConfig) error {
	server, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", config.host, config.port)),
		wish.WithHostKeyPath(config.hostKeyPath),
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

	log.Printf("Starting SSH server on %s:%d", config.host, config.port)
	go func() {
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, ssh.ErrServerClosed) {
			log.Printf("server error: %v", serveErr)
			done <- syscall.SIGTERM
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if shutdownErr := server.Shutdown(ctx); shutdownErr != nil && !errors.Is(shutdownErr, ssh.ErrServerClosed) {
		return shutdownErr
	}
	return nil
}

func teaHandler(_ ssh.Session) (tea.Model, []tea.ProgramOption) {
	model := asciimation.NewDefaultModel()
	model.AltScreen = true

	return model, nil
}
