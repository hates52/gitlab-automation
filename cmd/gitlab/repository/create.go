package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	gitlab "gitlab.com/gitlab-org/api/client-go"
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

	client, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	projectOptions := &gitlab.CreateProjectOptions{
		Name:        gitlab.Ptr(projectName),
		Path:        gitlab.Ptr(projectName),
		Description: gitlab.Ptr(projectDescription),
		NamespaceID: gitlab.Ptr(namespaceID),
		Visibility:  gitlab.Ptr(gitlab.VisibilityValue(visibility)),
	}

	project, _, err := client.Projects.CreateProject(projectOptions)
	if err != nil {
		log.Fatalf("Failed to create GitLab repository: %v", err)
	}

	fmt.Printf("Repository created successfully\n")
	fmt.Printf("Name: %s\n", project.Name)
	fmt.Printf("Description: %s\n", project.Description)
	fmt.Printf("Web URL: %s\n", project.WebURL)
}
