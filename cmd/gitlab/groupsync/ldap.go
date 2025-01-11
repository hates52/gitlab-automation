package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	common "github.com/Cloud-for-You/devops-cli/pkg"
	gitlab "github.com/Cloud-for-You/devops-cli/pkg/gitlab"
	ldap "github.com/Cloud-for-You/devops-cli/pkg/gitlab/groupsync/ldap"
	client "gitlab.com/gitlab-org/api/client-go"
)

var (
	ldapHost, ldapBindDN, ldapPassword, ldapSearchBase, ldapFilter string
)

var LdapCmd = &cobra.Command{
	Use:   "ldap",
	Short: "Synchronization Groups and Members from LDAP",
	Long: `The "groupsync" command allows you to synchronize groups and their members 
from an LDAP server to your GitLab instance. This is particularly useful for ensuring 
that group memberships are consistent and up-to-date, enabling efficient permissions 
management in GitLab.

The command connects to an LDAP server, retrieves group and user data, and updates the 
corresponding groups and members in GitLab. You can use various flags to specify the 
LDAP connection, the source groups to synchronize, and other options.

Examples:
  # Synchronize all groups from the default LDAP server
  devops-cli groupsync ldap \
	--ldapHost "ldaps://secure.example.com" \
	--ldapBindDN "CN=manager,DC=example,DC=com" \
	--ldapPassword "LDAP_Password_123" \
	--ldapSearchBase "OU=Groups,DC=example,DC=com" \
	--gitlabUrl "https://gitlab.example.com" \
	--gitlabToken "2fb5ae578dd22282da6289d1"
`,
	Run: ldapGroupSync,
}

func init() {
	// GitLab GroupSync LDAP
	LdapCmd.Flags().StringVarP(&ldapHost, "ldapHost", "H", "", "the IP address or resolvable name to use to connect to the directory server")
	viper.BindPFlag("ldapHost", LdapCmd.Flags().Lookup("ldapHost"))
	LdapCmd.Flags().StringVarP(&ldapBindDN, "ldapBindDN", "D", "", "the DN to use to bind to the directory server when performing simple authentication")
	viper.BindPFlag("ldapBindDN", LdapCmd.Flags().Lookup("ldapBindDN"))
	LdapCmd.Flags().StringVarP(&ldapPassword, "ldapPassword", "W", "", "the password to use to bind to the directory server when performing simple authentication or a password-based SASL mechanism")
	viper.BindPFlag("ldapPassword", LdapCmd.Flags().Lookup("ldapPassword"))
	LdapCmd.Flags().StringVarP(&ldapSearchBase, "ldapSearchBase", "b", "", "specifies the base DN that should be used for the search")
	viper.BindPFlag("ldapSearchBase", LdapCmd.Flags().Lookup("ldapSearchBase"))
	LdapCmd.Flags().StringVarP(&ldapFilter, "ldapGroupFilter", "f", "(objectClass=group)", "(optional) specified LDAP group search filter")
	viper.BindPFlag("ldapGroupFilter", LdapCmd.Flags().Lookup("ldapGroupFilter"))

	LdapCmd.MarkFlagRequired("ldapHost")
	LdapCmd.MarkFlagRequired("ldapBindDN")
	LdapCmd.MarkFlagRequired("ldapPassword")
	LdapCmd.MarkFlagRequired("ldapSearchBase")
}

func ldapGroupSync(cmd *cobra.Command, args []string) {
	ldapHost, _ := cmd.Flags().GetString("ldapHost")
	ldapBindDN, _ := cmd.Flags().GetString("ldapBindDN")
	ldapPassword, _ := cmd.Flags().GetString("ldapPassword")
	ldapSearchBase, _ := cmd.Flags().GetString("ldapSearchBase")
	ldapGroupFilter, _ := cmd.Flags().GetString("ldapGroupFilter")

	gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
	gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

	// Overeni, ze mame gitlabURL a gitlabToken
	if gitlabToken == "" || gitlabUrl == "" {
		log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
	}

	ldapConfig := ldap.LDAPConfig{
		Host:     ldapHost,
		BindDN:   ldapBindDN,
		Password: ldapPassword,
		BaseDN:   ldapSearchBase,
	}
	connector, err := ldap.NewLDAPConnector(ldapConfig)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	defer connector.Close()

	groupSyncer := ldap.NewLDAPGroupSyncer(connector, ldapGroupFilter)

	// Nacteni skupin z LDAPu
	groups, err := groupSyncer.GetLdapGroups()
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	client, err := client.NewClient(gitlabToken, client.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	// Iterace pres vsechny LDAP skupiny
	for _, group := range groups.Entries {
		groupName := group.GetAttributeValue("cn")

		if gitlabToken == "" || gitlabUrl == "" {
			log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
		}

		// Ziskani seznamu clenu skupiny z LDAPu
		ldapMembersDNs, err := groupSyncer.ListLdapGroupMemberDNs(group.DN)
		if err != nil {
			log.Fatalf("Error listing Ldap group members: %v", err)
		}

		var ldapMembers []common.Member
		for _, dn := range ldapMembersDNs {
			// Ziskani atributu clena
			//attributes := []string{"givenName"}
			userAttributes, err := connector.GetLdapUserAttributes(dn, nil)
			if err != nil {
				log.Fatalf("Error to get attributes for user %s: %v", dn, err)
			}
			if len(userAttributes.Entries) > 0 {
				mail := userAttributes.Entries[0].GetAttributeValue("sAMAccountName")
				ldapMembers = append(ldapMembers, common.Member{Name: mail})
			}
		}

		// Ziskani clenu skupiny z GitLab
		gitlabMembersRaw, err := gitlab.ListGitlabGroupMembers(client, groupName)
		if err != nil {
			log.Fatalf("Error listing GitLab group members: %v", err)
		}

		var gitlabMembers []common.Member
		for _, member := range gitlabMembersRaw {
			if member.Username != "root" {
				gitlabMembers = append(gitlabMembers, common.Member{Name: member.Username})
			}
		}

		// Overime zda skupina podle nazvu v GitLabu existuje a pokud ano, pouze ji sesynchronizujeme
		// a pokracujeme dalsi skupinou
		group, err := gitlab.GetGroup(client, groupName)
		if err == nil {
			fmt.Printf("Synchronizing members of an existing GitLab group [%s]\n", group.Name)

			missing, extra := common.CompareMembers(gitlabMembers, ldapMembers)
			for _, m := range missing {
				fmt.Printf("Add members %s to GitLab group %s\n", m.Name, group.Name)
				err = gitlab.AddMemberToGroup(client, group.Name, m.Name)
				if err != nil {
					log.Fatalf("error: %v", err)
				} else {
					fmt.Printf("User '%s' successfully added to the group '%s'.", m.Name, group.Name)
				}
			}
			for _, m := range extra {
				fmt.Printf("Remove member %s from GitLab group %s\n", m.Name, group.Name)
				err = gitlab.RemoveUserFromGroup(client, group.Name, m.Name)
				if err != nil {
					log.Fatalf("error: %v", err)
				} else {
					fmt.Printf("User '%s' successfully remove from group '%s'.", m.Name, group.Name)
				}
			}

			continue
		}

		// Zalozime skupinu v GitLabu a vlozime do ni membery
		_, response, err := gitlab.CreateGroup(client, groupName, "", "private")
		if err != nil {
			if response != nil && response.StatusCode == http.StatusConflict {
				fmt.Printf("Group '%s' is exists.\n", groupName)
			} else {
				fmt.Printf("Failed to create GitLab group '%s': %v\n", groupName, err)
			}
		}
	}
}
