package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	ldap, err := ldap.NewLDAPGroupSyncer(ldapHost, ldapBindDN, ldapPassword, ldapSearchBase, ldapFilter)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	defer ldap.Close()

	groups, err := ldap.GetLdapGroups()
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
	gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

	if gitlabToken == "" || gitlabUrl == "" {
		log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
	}

	client, err := client.NewClient(gitlabToken, client.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}
	
	for _, group := range groups.Entries {
		groupName := group.GetAttributeValue("cn")
		gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
		gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

		if gitlabToken == "" || gitlabUrl == "" {
			log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
		}

    // Ziskame seznam uzivatelu, kteri maji byt v dane skupine nastaveni
		ldapMembers, err := ldap.ListLdapGroupMembers(groupName)
		if err != nil {
      fmt.Printf("Chyba v ziskavani clenu skupiny %s z LDAPu.\n", groupName)
		}
		fmt.Println(ldapMembers)

    
    gitlabMembers, err := gitlab.ListGitlabGroupMembers(client, groupName)
		if err != nil {
      fmt.Printf("Chyba v ziskavani clenu skupiny %s z LDAPu.\n", groupName)
		}
		fmt.Println(gitlabMembers)

		// Overime zda skupina podle nazvu v GitLabu existuje a pokud ano, pouze do ni membery vlozime
		// a pokracujeme dalsi skupinou
		result, err := gitlab.GetGroup(client, groupName)
		if err == nil {
			fmt.Printf("Group %v exist.\n", result.Name)
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

		// Pro skupinu ziskame membery a vlozime je do skupiny
		//members, err := groupsync.ListLdapGroupMembers(group.DN)
		//if err != nil {
		//	log.Fatalf("ERROR: %s", err)
		//}
		//for _, member := range members {
		//	fmt.Printf("  - %s\n", member)
		//}
	}
}
