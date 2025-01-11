/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	group "github.com/Cloud-for-You/devops-cli/cmd/gitlab/group"
	groupsync "github.com/Cloud-for-You/devops-cli/cmd/gitlab/groupsync"
	list "github.com/Cloud-for-You/devops-cli/cmd/gitlab/list"
	project "github.com/Cloud-for-You/devops-cli/cmd/gitlab/project"
	gitlab "github.com/Cloud-for-You/devops-cli/pkg/gitlab"
	client "gitlab.com/gitlab-org/api/client-go"
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

// Get ma
var WhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display the username currently user",
	Run:   whoami,
}

func init() {
	GitlabCmd.AddCommand(WhoamiCmd)
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

func whoami(cmd *cobra.Command, args []string) {
	gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
	gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

	if gitlabToken == "" || gitlabUrl == "" {
		log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
	}

	client, err := client.NewClient(gitlabToken, client.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	username, _ := gitlab.Whoami(client)
	fmt.Println(*username)

}
