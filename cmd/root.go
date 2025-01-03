/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	gitlab "github.com/Cloud-for-You/devops-cli/cmd/gitlab"
)

var (
	Debug      bool
	configFile string
	Command    string
	Flags      map[string]string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                   "devops-cli",
	Short:                 "Client for DEVOPS tools management",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		if Command != "" {
			executeCommandFromConfig(cmd)
		} else {
			fmt.Println("Run `devops-cli --help` for usage.")
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(gitlab.GitlabCmd)

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	rootCmd.Flags().StringVar(&configFile, "config", "", "Configuration file from command and flags.")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		//fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	Command = viper.GetString("command")
	Flags = viper.GetStringMapString("flags")
	fmt.Println(Command)
}

func executeCommandFromConfig(cmd *cobra.Command) {
	parts := strings.Split(Command, "/")

	// Pripravime argumenty pro prikaz
	args := append([]string{parts[0]}, parts[1:]...)

	for flag, value := range Flags {
		args = append(args, fmt.Sprintf("--%s=\"%s\"", flag, value))
	}

	// Nastavime argumenty prikazu
	cmd.SetArgs(args)
	fmt.Printf("Executing command: %s\n", strings.Join(args, " "))

	// Spustime prikaz s argumenty
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
}
