/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	gitlabUrl, gitlabToken string
  ldapHost, ldapBindDN, ldapPassword, ldapSearchBase string
  azureTenantID string
)

var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "Managing GitLab repository",
	DisableFlagsInUseLine: true,
}

var groupSyncCmd = &cobra.Command{
	Use:   "groupsync",
	Short: "Synchronization Groups and Members to GitLab",
	DisableFlagsInUseLine: true,
}

var ldapGroupSyncCmd = &cobra.Command{
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
	Run: func(cmd *cobra.Command, args []string) {
		if Debug {
			for key, value := range viper.GetViper().AllSettings() {
				log.WithFields(log.Fields{
					key: value,
				}).Info("Command Flag")
			}
		}

		fmt.Println("LDAP groupsync",)
	},
}

var azureGroupSyncCmd = &cobra.Command{
	Use:   "azure",
	Short: "Synchronization Groups and Members from Azure EntryID",
	Run: func(cmd *cobra.Command, args []string) {
		if Debug {
			for key, value := range viper.GetViper().AllSettings() {
				log.WithFields(log.Fields{
					key: value,
				}).Info("Command Flag")
			}
		}

		fmt.Println("Azure EntryID groupsync",)
	},
}

var fileGroupSyncCmd = &cobra.Command{
	Use:   "file",
	Short: "Synchronization Groups and Members from file",
	Run: func(cmd *cobra.Command, args []string) {
		if Debug {
			for key, value := range viper.GetViper().AllSettings() {
				log.WithFields(log.Fields{
					key: value,
				}).Info("Command Flag")
			}
		}

		fmt.Println("File groupsync",)
	},
}

func init() {
	rootCmd.AddCommand(gitlabCmd)
	gitlabCmd.AddCommand(groupSyncCmd)
  groupSyncCmd.AddCommand(ldapGroupSyncCmd)
	groupSyncCmd.AddCommand(azureGroupSyncCmd)
	groupSyncCmd.AddCommand(fileGroupSyncCmd)

  // Perzistentni flagy pro GITLAB
	gitlabCmd.PersistentFlags().StringVar(&gitlabUrl, "gitlabUrl", "", "GitLab URL adresses")
	viper.BindPFlag("gitlabUrl", gitlabCmd.PersistentFlags().Lookup("gitlabUrl"))
	gitlabCmd.PersistentFlags().StringVar(&gitlabToken, "gitlabToken", "", "login token")
	viper.BindPFlag("gitlabToken", gitlabCmd.PersistentFlags().Lookup("gitlabToken"))

  // Flagy pro LDAP
	ldapGroupSyncCmd.Flags().StringVarP(&ldapHost, "ldapHost", "H", "", "the IP address or resolvable name to use to connect to the directory server")
	ldapGroupSyncCmd.MarkFlagRequired("ldapHost")
	viper.BindPFlag("ldapHost", ldapGroupSyncCmd.Flags().Lookup("ldapHost"))
	ldapGroupSyncCmd.Flags().StringVarP(&ldapBindDN, "ldapBindDN", "D", "", "the DN to use to bind to the directory server when performing simple authentication")
	ldapGroupSyncCmd.MarkFlagRequired("ldapBindDN")
	viper.BindPFlag("ldapBindDN", ldapGroupSyncCmd.Flags().Lookup("ldapBindDN"))
	ldapGroupSyncCmd.Flags().StringVarP(&ldapPassword, "ldapPassword", "W", "", "the password to use to bind to the directory server when performing simple authentication or a password-based SASL mechanism")
	ldapGroupSyncCmd.MarkFlagRequired("ldapPassword")
	viper.BindPFlag("ldapPassword", ldapGroupSyncCmd.Flags().Lookup("ldapPassword"))
	ldapGroupSyncCmd.Flags().StringVarP(&ldapSearchBase, "ldapSearchBase", "b", "", "specifies the base DN that should be used for the search")
	ldapGroupSyncCmd.MarkFlagRequired("ldapSearchBase")
	viper.BindPFlag("ldapSearchBase", ldapGroupSyncCmd.Flags().Lookup("ldapSearchBase"))

	// Flagy pro Azure
	azureGroupSyncCmd.Flags().StringVar(&azureTenantID, "tenantID", "", "Tenant ID (default \"\")")
	viper.BindPFlag("tenantID", azureGroupSyncCmd.Flags().Lookup("tenantID"))
}
