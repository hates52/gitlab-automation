package cmd

import "github.com/spf13/cobra"

var GroupCmd = &cobra.Command{
	Use:                   "group",
	Short:                 "Managing GitLab groups",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	GroupCmd.AddCommand(CreateCmd)
}
