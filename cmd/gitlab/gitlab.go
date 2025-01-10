/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	group "github.com/Cloud-for-You/devops-cli/cmd/gitlab/group"
	groupsync "github.com/Cloud-for-You/devops-cli/cmd/gitlab/groupsync"
	list "github.com/Cloud-for-You/devops-cli/cmd/gitlab/list"
	project "github.com/Cloud-for-You/devops-cli/cmd/gitlab/project"
)

var (
	gitlabUrl, gitlabToken string
)

var GitlabCmd = &cobra.Command{
	Use:                   "gitlab",
	Short:                 "Managing GitLab platform",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	GitlabCmd.AddCommand(list.ListCmd)
	GitlabCmd.AddCommand(project.RepositoryCmd)
	GitlabCmd.AddCommand(group.GroupCmd)
	GitlabCmd.AddCommand(groupsync.GroupSyncCmd)

	// FLAGS
	GitlabCmd.PersistentFlags().StringVar(&gitlabUrl, "gitlabUrl", "", "GitLab URL adresses")
	viper.BindPFlag("gitlabUrl", GitlabCmd.PersistentFlags().Lookup("gitlabUrl"))
	GitlabCmd.PersistentFlags().StringVar(&gitlabToken, "gitlabToken", "", "login token")
	viper.BindPFlag("gitlabToken", GitlabCmd.PersistentFlags().Lookup("gitlabToken"))

	GitlabCmd.MarkPersistentFlagRequired("gitlabUrl")
	GitlabCmd.MarkPersistentFlagRequired("gitlabToken")
}
