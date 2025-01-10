package common

import "fmt"

type Member struct {
	Name string
}

// currentMembers is members in GitLab group
// desiredMembers is members in SRC group (LDAP, Azure, File, etc...)
// missing: User is LDAP but not in GitLab (Create)
// extra: User in GitLab but not in LDAP (Delete)
func CompareMembers(gitlabMembers, sourceMembers []Member) (missing, extra []Member) {
	// Vytvorime mapy pro rychle vyhledavani
	gitlabSet := make(map[string]struct{})
	srcSet := make(map[string]struct{})

	// Naplnime mapu
	for _, m := range gitlabMembers {
		gitlabSet[m.Name] = struct{}{}
	}
	for _, m := range sourceMembers {
		srcSet[m.Name] = struct{}{}
	}

	// Najdeme chybejici cleny (v source, ale ne v GitLab)
	for _, m := range sourceMembers {
		if _, exists := gitlabSet[m.Name]; !exists {
			missing = append(missing, m)
		}
	}
	fmt.Printf("Created: %v\n", missing)

	// Najdeme prebyvajici cleny (v GitLab, ale ne v source)
	for _, m := range gitlabMembers {
		if _, exists := srcSet[m.Name]; !exists {
			extra = append(extra, m)
		}
	}
	fmt.Printf("Deleted: %v\n", extra)

	return missing, extra
}
