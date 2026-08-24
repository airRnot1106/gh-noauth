package authstub

import (
	"bytes"
	"testing"

	"github.com/cli/cli/v2/pkg/cmdutil"
)

func TestNewCmdAuth_subcommandsDisabled(t *testing.T) {
	tests := []string{
		"login",
		"logout",
		"status",
		"refresh",
		"git-credential",
		"setup-git",
		"token",
		"switch",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := NewCmdAuth(&cmdutil.Factory{})
			cmd.SetArgs([]string{name, "--some-unknown-flag", "positional"})

			out := &bytes.Buffer{}
			cmd.SetOut(out)
			cmd.SetErr(out)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("expected an error, got none")
			}
			if err.Error() != disabledMessage {
				t.Errorf("got error %q, want %q", err.Error(), disabledMessage)
			}
		})
	}
}

func TestNewCmdAuth_hidden(t *testing.T) {
	cmd := NewCmdAuth(&cmdutil.Factory{})

	if !cmd.Hidden {
		t.Error("expected auth command to be hidden")
	}
	for _, sub := range cmd.Commands() {
		if !sub.Hidden {
			t.Errorf("expected subcommand %q to be hidden", sub.Name())
		}
	}
}

func TestNewCmdAuth_parentAlone(t *testing.T) {
	cmd := NewCmdAuth(&cmdutil.Factory{})
	cmd.SetArgs([]string{})

	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)

	err := cmd.Execute()
	if err == nil || err.Error() != disabledMessage {
		t.Errorf("got error %v, want %q", err, disabledMessage)
	}
}
