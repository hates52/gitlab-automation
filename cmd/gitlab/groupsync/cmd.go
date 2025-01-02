package cmd

import "github.com/spf13/cobra"

var GroupSyncCmd = &cobra.Command{
	Use:                   "groupsync",
	Short:                 "Synchronization Groups and Members to GitLab",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	GroupSyncCmd.AddCommand(LdapCmd)
}
