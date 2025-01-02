package devops_cli

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func GetGroups(client *gitlab.Client) ([]*gitlab.Group, error) {
	var allGroups []*gitlab.Group
	page := 1
	perPage := 20

	for {
		options := &gitlab.ListGroupsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		}

		groups, res, err := client.Groups.ListGroups(options)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve groups: %v", err)
		}

		allGroups = append(allGroups, groups...)

		if res.CurrentPage >= res.TotalPages {
			break
		}

		page++
	}

	return allGroups, nil
}

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
