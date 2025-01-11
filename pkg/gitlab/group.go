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

func ListGitlabGroupMembers(client *gitlab.Client, groupId string) ([]*gitlab.GroupMember, error) {
	var allMembers []*gitlab.GroupMember
	page := 1
	perPage := 20

	for {
		options := &gitlab.ListGroupMembersOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		}

		groups, res, err := client.Groups.ListGroupMembers(groupId, options)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve groupMembers: %v", err)
		}

		allMembers = append(allMembers, groups...)

		if res.CurrentPage >= res.TotalPages {
			break
		}

		page++
	}

	return allMembers, nil
}

func GetGroup(client *gitlab.Client, groupName string) (*gitlab.Group, error) {
	group, _, err := client.Groups.GetGroup(groupName, nil)
	if err != nil {
		return nil, err
	}
	return group, nil
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

// AddUserToGroup adds a user to a group
func AddMemberToGroup(client *gitlab.Client, groupname string, username string) error {
	// Na zaklade jmena skupiny ziskame jeji ID
	groups, _, err := client.Groups.ListGroups(&gitlab.ListGroupsOptions{
		Search: &groupname,
	})
	if err != nil {
		return fmt.Errorf("error retrieving group: %w", err)
	}

	var groupID int
	for _, group := range groups {
		if group.Name == groupname {
			groupID = group.ID
			break
		}
	}

	if groupID == 0 {
		return fmt.Errorf("group '%s' not found", groupname)
	}

	// Na zaklade jmena uzivatele ziskame jeho ID
	users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{
		Username: &username,
	})
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}

	var userID int
	for _, user := range users {
		if user.Name == username {
			userID = user.ID
			break
		}
	}

	if userID == 0 {
		return fmt.Errorf("user '%s' not found", username)
	}

	// Pridame uzivatele do skupiny
	_, _, err = client.GroupMembers.AddGroupMember(groupID, &gitlab.AddGroupMemberOptions{
		UserID:      &userID,
		AccessLevel: gitlab.Ptr(gitlab.DeveloperPermissions),
	})
	if err != nil {
		return fmt.Errorf("error adding user to group: %w", err)
	}

	return nil
}

// RemoveUserFromGroup removes a user from a group
func RemoveUserFromGroup(client *gitlab.Client, groupname string, username string) error {
	// Retrieve group ID by name
	groups, _, err := client.Groups.ListGroups(&gitlab.ListGroupsOptions{
		Search: &groupname,
	})
	if err != nil {
		return fmt.Errorf("error retrieving group: %w", err)
	}

	var groupID int
	for _, group := range groups {
		if group.Name == groupname {
			groupID = group.ID
			break
		}
	}

	if groupID == 0 {
		return fmt.Errorf("group '%s' not found", groupname)
	}

	// Retrieve user ID by username
	users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{
		Username: &username,
	})
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}

	var userID int
	for _, user := range users {
		if user.Username == username {
			userID = user.ID
			break
		}
	}

	if userID == 0 {
		return fmt.Errorf("user '%s' not found", username)
	}

	// Remove user from group
	_, err = client.GroupMembers.RemoveGroupMember(groupID, userID, nil)
	if err != nil {
		return fmt.Errorf("error removing user from group: %w", err)
	}

	return nil
}
