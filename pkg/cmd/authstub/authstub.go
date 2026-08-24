// Package authstub is wired into pkg/cmd/root in place of pkg/cmd/auth.
// It disables `gh auth` and every one of its subcommands so that
// authentication comes exclusively from the GH_TOKEN/GITHUB_TOKEN
// environment variables or a pre-populated hosts.yml, both of which
// keep working unaffected.
package authstub

import (
	"errors"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

const disabledMessage = "Auth command has been disabled."

var disabledSubcommands = []string{
	"login",
	"logout",
	"status",
	"refresh",
	"git-credential",
	"setup-git",
	"token",
	"switch",
}

// NewCmdAuth returns a stub `gh auth` command tree: every
// `gh auth <subcommand>` invocation fails with a single consistent
// message instead of running real authentication logic or falling
// through to an "unknown command" error.
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "auth <command>",
		Short:              "Authenticate gh and git with GitHub",
		GroupID:            "core",
		Hidden:             true,
		DisableFlagParsing: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
		RunE:               runDisabled,
	}

	// Bypasses the root command's not-authenticated check so invocations
	// show the disabled message instead of a generic auth prompt.
	cmdutil.DisableAuthCheck(cmd)

	for _, name := range disabledSubcommands {
		cmd.AddCommand(&cobra.Command{
			Use:                name,
			Hidden:             true,
			DisableFlagParsing: true,
			SilenceUsage:       true,
			SilenceErrors:      true,
			RunE:               runDisabled,
		})
	}

	cmdutil.DisableTelemetryForSubcommands(cmd)

	return cmd
}

func runDisabled(*cobra.Command, []string) error {
	return errors.New(disabledMessage)
}
