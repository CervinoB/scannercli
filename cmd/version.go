/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
	"go.k6.io/k6/cmd/state"
	"go.k6.io/k6/lib/consts"
)

const (
	commitKey      = "commit"
	commitDirtyKey = "commit_dirty"
)

func versionDetails() map[string]interface{} {
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

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return details
	}
	var (
		commit string
		dirty  bool
	)
	for _, s := range buildInfo.Settings {
		switch s.Key {
		case "vcs.revision":
			commitLen := 10
			if len(s.Value) < commitLen {
				commitLen = len(s.Value)
			}
			commit = s.Value[:commitLen]
		case "vcs.modified":
			if s.Value == "true" {
				dirty = true
			}
		default:
		}
	}
	if commit == "" {
		return details
	}

	details[commitKey] = commit
	if dirty {
		details[commitDirtyKey] = true
	}

	return details
}

type versionCmd struct {
	isJSON bool
}

func (c *versionCmd) run(cmd *cobra.Command, _ []string) error {
	if !c.isJSON {
		root := cmd.Root()
		root.SetArgs([]string{"--version"})
		_ = root.Execute()
		return nil
	}

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

// versionCmd represents the version command
var versionCmdOld = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version called")
	},
}

func init() {
	rootCmd.AddCommand(versionCmdOld)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
