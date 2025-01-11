package cmd

import (
	"github.com/spf13/cobra"
)

var RepositoryCmd = &cobra.Command{
	Use:                   "project",
	Short:                 "Managing GitLab project",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RepositoryCmd.AddCommand(CreateCmd)
}
