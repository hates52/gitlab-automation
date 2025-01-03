package gitlab

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func ListGroups(client *gitlab.Client) ([]*gitlab.Group, error) {
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

func CreateGroup(client *gitlab.Client, groupName string, groupDescription string, visibility string) (*gitlab.Group, *gitlab.Response, error) {

	groupOptions := &gitlab.CreateGroupOptions{
		Name:        gitlab.Ptr(groupName),
		Path:        gitlab.Ptr(groupName),
		Description: gitlab.Ptr(groupDescription),
		Visibility:  gitlab.Ptr(gitlab.VisibilityValue(visibility)),
	}

	group, res, err := client.Groups.CreateGroup(groupOptions)
	if err != nil {
		return nil, res, err
	}

	return group, res, nil
}
