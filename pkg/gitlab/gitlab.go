package gitlab

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Whoami(client *gitlab.Client) (*string, error) {
	user, _, err := client.Users.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("error retrieving current user: %w", err)
	}
	return &user.Username, nil
}
