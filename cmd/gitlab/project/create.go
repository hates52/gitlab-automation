package cmd

import (
	"fmt"
	"log"
	"net/http"

	gitlab "github.com/Cloud-for-You/devops-cli/pkg/gitlab"
	"github.com/spf13/cobra"
	client "gitlab.com/gitlab-org/api/client-go"
)

var (
	projectName        string
	projectDescription string
	namespaceID        int
	visibility         string
)

// Create GitLab repository
var CreateCmd = &cobra.Command{
	Use:                   "create",
	Short:                 "Create GitLab repository",
	DisableFlagsInUseLine: true,
	Run:                   createRepository,
}

func init() {
	CreateCmd.Flags().StringVar(&projectName, "name", "", "Name of the repository (required)")
	CreateCmd.Flags().StringVar(&projectDescription, "description", "", "Description of the repository")
	CreateCmd.Flags().IntVar(&namespaceID, "namespace", 0, "Namespace ID under which the repository will be created")
	CreateCmd.Flags().StringVar(&visibility, "visibility", "private", "Visibility of the repository (private, internal, public)")

	CreateCmd.MarkFlagRequired("name")
}

func createRepository(cmd *cobra.Command, args []string) {
	gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
	gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

	if gitlabToken == "" || gitlabUrl == "" {
		log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
	}

	client, err := client.NewClient(gitlabToken, client.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	result, res, err := gitlab.CreateProject(client, projectName, namespaceID, projectDescription, visibility)
	if err != nil {
		if res != nil && res.StatusCode == http.StatusConflict {
			fmt.Printf("Project '%s' is exists.\n", projectName)
		} else {
			fmt.Printf("Failed to create GitLab project '%s': %v\n", projectName, err)
		}
	} else {
		fmt.Printf("Project created successfully\n")
		fmt.Printf("Name: %s\n", result.Name)
		fmt.Printf("Description: %s\n", result.Description)
		fmt.Printf("Web URL: %s\n", result.WebURL)
	}
}
