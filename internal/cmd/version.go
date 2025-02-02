/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/CervinoB/scannercli/cmd/state"
	"github.com/CervinoB/scannercli/lib/consts"

	"github.com/spf13/cobra"
)

// fullVersion returns the maximally full version and build information for
// the currently running k6 executable.
func fullVersion() string {
	details := getVersion()

	goVersionArch := fmt.Sprintf("%s, %s/%s", details["go_version"], details["go_os"], details["go_arch"])

	cliversion := fmt.Sprintf("%s", details["version"])
	cliversion = strings.TrimLeft(cliversion, "v")

	return fmt.Sprintf("%s (%s)", cliversion, goVersionArch)
}

func getVersion() map[string]interface{} {
	v := consts.Version
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}

	details := map[string]interface{}{
		"version":    v,
		"go_version": runtime.Version(),
		"go_os":      runtime.GOOS,
		"go_arch":    runtime.GOARCH,
	}

	return details
}

func versionString() string {
	v := fullVersion()
	return v
}

type versionCmd struct {
	gs     *state.GlobalState
	isJSON bool
}

func (c *versionCmd) run(cmd *cobra.Command, _ []string) error {
	if !c.isJSON {
		root := cmd.Root()
		root.SetArgs([]string{"--version"})
		_ = root.Execute()
		return nil
	}
	details := versionDetails()
	fmt.Println(details["version"])    // version
	fmt.Println(details["go_version"]) // go version
	fmt.Println(details["go_os"])      // go os
	fmt.Println(details["go_arch"])    // go arch
	return nil
}

func getCmdVersion(gs *state.GlobalState) *cobra.Command {
	versionCmd := &versionCmd{gs: gs}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show application version",
		Long:  `Show the application version and exit.`,
		RunE:  versionCmd.run,
	}

	cmd.Flags().BoolVar(&versionCmd.isJSON, "json", false, "if set, output version information will be in JSON format")

	return cmd
}
