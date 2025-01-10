package cmd

import (
	"github.com/spf13/cobra"
)

var RepositoryCmd = &cobra.Command{
	Use:                   "repository",
	Short:                 "Managing GitLab repository",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RepositoryCmd.AddCommand(CreateCmd)
}
