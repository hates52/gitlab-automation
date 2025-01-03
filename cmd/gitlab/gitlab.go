/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	group "github.com/Cloud-for-You/devops-cli/cmd/gitlab/group"
	groupsync "github.com/Cloud-for-You/devops-cli/cmd/gitlab/groupsync"
	repository "github.com/Cloud-for-You/devops-cli/cmd/gitlab/repository"

	devops_cli "github.com/Cloud-for-You/devops-cli/pkg/gitlab"
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

// Get all Project name and ID from Gitlab
var ListProjectsCmd = &cobra.Command{
	Use:   "listProjects",
	Short: "Get all GitLab project",
	Run:   listProjects,
}

// Get all Project name and ID from Gitlab
var ListGroupsCmd = &cobra.Command{
	Use:   "listGroups",
	Short: "List all GitLab groups",
	Run:   listGroups,
}

func init() {
	GitlabCmd.AddCommand(repository.RepositoryCmd)
	GitlabCmd.AddCommand(group.GroupCmd)
	GitlabCmd.AddCommand(groupsync.GroupSyncCmd)

	GitlabCmd.AddCommand(ListProjectsCmd)
	GitlabCmd.AddCommand(ListGroupsCmd)

	// FLAGS
	GitlabCmd.PersistentFlags().StringVar(&gitlabUrl, "gitlabUrl", "", "GitLab URL adresses")
	viper.BindPFlag("gitlabUrl", GitlabCmd.PersistentFlags().Lookup("gitlabUrl"))
	GitlabCmd.PersistentFlags().StringVar(&gitlabToken, "gitlabToken", "", "login token")
	viper.BindPFlag("gitlabToken", GitlabCmd.PersistentFlags().Lookup("gitlabToken"))

	GitlabCmd.MarkPersistentFlagRequired("gitlabUrl")
	GitlabCmd.MarkPersistentFlagRequired("gitlabToken")
}

func listProjects(cmd *cobra.Command, args []string) {
	gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
	gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

	if gitlabToken == "" || gitlabUrl == "" {
		log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
	}

	client, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	projects, err := devops_cli.GetProjects(client)
	if err != nil {
		log.Fatalf("Error retrieving projects :%v", err)
	}

	for _, project := range projects {
		fmt.Printf("ID: %d, Name: %s\n", project.ID, project.Name)
	}
}

func listGroups(cmd *cobra.Command, args []string) {
	gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
	gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

	if gitlabToken == "" || gitlabUrl == "" {
		log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
	}

	client, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	groups, err := devops_cli.GetGroups(client)
	if err != nil {
		log.Fatalf("Error retrieving projects :%v", err)
	}

	for _, group := range groups {
		fmt.Printf("ID: %d, Name: %s\n", group.ID, group.Name)
	}
}
