package cmd

import (
	"fmt"
	"log"

	gitlab "github.com/Cloud-for-You/devops-cli/pkg/gitlab"
	"github.com/spf13/cobra"
	client "gitlab.com/gitlab-org/api/client-go"
)

var ListCmd = &cobra.Command{
	Use:                   "list",
	Short:                 "List GitLab object",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Get all Project name and ID from Gitlab
var ListProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Get all GitLab project",
	Run:   listProjects,
}

// Get all Project name and ID from Gitlab
var ListGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all GitLab groups",
	Run:   listGroups,
}

func init() {
	ListCmd.AddCommand(ListProjectsCmd)
	ListCmd.AddCommand(ListGroupsCmd)
}

func listProjects(cmd *cobra.Command, args []string) {
	gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
	gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

	if gitlabToken == "" || gitlabUrl == "" {
		log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
	}

	client, err := client.NewClient(gitlabToken, client.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	projects, err := gitlab.GetProjects(client)
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

	client, err := client.NewClient(gitlabToken, client.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	groups, err := gitlab.ListGroups(client)
	if err != nil {
		log.Fatalf("Error retrieving projects :%v", err)
	}

	for _, group := range groups {
		fmt.Printf("ID: %d, Name: %s\n", group.ID, group.Name)
	}
}
