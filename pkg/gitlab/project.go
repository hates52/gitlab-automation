package gitlab

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func GetProjects(client *gitlab.Client) ([]*gitlab.Project, error) {
	var allProjects []*gitlab.Project
	page := 1
	perPage := 20

	for {
		options := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		}

		projects, res, err := client.Projects.ListProjects(options)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve projects: %v", err)
		}

		allProjects = append(allProjects, projects...)

		if res.CurrentPage >= res.TotalPages {
			break
		}

		page++
	}

	return allProjects, nil
}
