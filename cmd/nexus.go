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

// nexusCmd represents the nexus command
var nexusCmd = &cobra.Command{
	Use:   "nexus",
	Short: "Managing Nexus registry",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		if Debug {
			for key, value := range viper.GetViper().AllSettings() {
				log.WithFields(log.Fields{
					key: value,
				}).Info("Command Flag")
			}
		}
		fmt.Println("nexus called")
	},
}

func init() {
	rootCmd.AddCommand(nexusCmd)
}
