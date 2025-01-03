package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Cloud-for-You/devops-cli/pkg/gitlab/groupsync"
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
	LdapCmd.Flags().StringVarP(&ldapFilter, "ldapGroupFilter", "f", "(CN=*)", "(optional) specified LDAP group search filter")
	viper.BindPFlag("ldapGroupFilter", LdapCmd.Flags().Lookup("ldapGroupFilter"))

	LdapCmd.MarkFlagRequired("ldapHost")
	LdapCmd.MarkFlagRequired("ldapBindDN")
	LdapCmd.MarkFlagRequired("ldapPassword")
	LdapCmd.MarkFlagRequired("ldapSearchBase")
}

func ldapGroupSync(cmd *cobra.Command, args []string) {
	client, err := groupsync.NewLDAPGroupSyncer(ldapHost, ldapBindDN, ldapPassword, ldapSearchBase, ldapFilter)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	defer client.Close()

	groups, err := client.GetGroups()
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	for _, group := range groups.Entries {

		fmt.Println(group.DN)
		members, err := client.GetMembers(group.DN)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}

		for _, member := range members {
			fmt.Printf("  - %s\n", member)
		}
	}

}
