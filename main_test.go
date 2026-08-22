package main

import (
	"strconv"
	"testing"
)

func TestNewRootCommandDefaults(t *testing.T) {
	command := newRootCommand()

	tests := []struct {
		name string
		want string
	}{
		{name: "host", want: defaultHost},
		{name: "port", want: strconv.Itoa(defaultPort)},
		{name: "host-key-path", want: defaultHostKeyPath},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flag := command.Flags().Lookup(test.name)
			if flag == nil {
				t.Fatalf("expected %q flag to exist", test.name)
			}
			if flag.DefValue != test.want {
				t.Fatalf("expected %q default %q, got %q", test.name, test.want, flag.DefValue)
			}
		})
	}
}
