/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/CervinoB/sonarcli/cmd/state"
	"github.com/spf13/cobra"
)

// Panic if the given error is not nil.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckIfError(gs *state.GlobalState, err error) {
	if err == nil {
		return
	}
	gs.Logger.Errorf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func exactArgsWithMsg(n int, msg string) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %d arg(s), received %d: %s", n, len(args), msg)
		}
		return nil
	}
}

func printToStdout(gs *state.GlobalState, s string) {
	if _, err := fmt.Fprint(os.Stdout, s); err != nil {
		gs.Logger.Errorf("could not print '%s' to stdout: %s", s, err.Error())
	}
}
