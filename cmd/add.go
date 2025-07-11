/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/CervinoB/scannercli/internal/todo"
	"github.com/spf13/cobra"
)

var priority int

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Deprecated: "This command is deprecated and will be removed in future versions. Used only to test and setup project.",
	Run:        addRun,
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().IntVarP(&priority, "priority", "p", 2, "Priority of the item (1: low, 2: medium, 3: high)")
}

func addRun(cmd *cobra.Command, args []string) {
	fmt.Println("add called")

	items, err := todo.ReadItems(dataFile)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read items: %v", err))
	}

	for _, x := range args {
		item := todo.Item{Text: x}
		item.SetPriority(priority)
		items = append(items, item)
	}
	err = todo.SaveItems(dataFile, items)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to save items: %v", err))
	}
	// fmt.Printf("%#v\n", items)

}
