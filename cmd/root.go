// Package cmd is the command line interface for the cache_clone application
package cmd

import (
	"os"

	"github.com/natemarks/cache_clone/config"
	"github.com/spf13/cobra"
)

var verbose bool
var mirror string

//var mirror, local, remote, secretID, userKey, tokenKey string

var settings config.Settings

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cache_clone",
	Short: "Clone a repo using a local mirror",
	Long: `Using credentials stored in AWS Secret Manager, create/update a local mirror.
Then use the local mirror to create local clone`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&settings.Verbose, "verbose", "v", false, "enable verbose/debug logging")

	rootCmd.PersistentFlags().StringVarP(&settings.Mirror, "mirror", "m", "", "Root location for all mirror repos")
	rootCmd.MarkFlagRequired("mirror")

	rootCmd.PersistentFlags().StringVarP(&settings.Local, "local", "l", "", "Location to create the repo clone")
	rootCmd.MarkFlagRequired("local")

	rootCmd.PersistentFlags().StringVarP(&settings.Remote, "remote", "r", "", "git remote server url. example: https://my.git.com/my/project.git")
	rootCmd.MarkFlagRequired("remote")

	rootCmd.PersistentFlags().StringVarP(&settings.SecretID, "secretID", "s", "", "AWS Secret Manager secretID path")
	rootCmd.MarkFlagRequired("secretID")

	rootCmd.PersistentFlags().StringVarP(&settings.UserKey, "userKey", "u", "", "username key in the secret JSON dict")
	rootCmd.MarkFlagRequired("userKey")

	rootCmd.PersistentFlags().StringVarP(&settings.TokenKey, "tokenKey", "t", "", "token key in the secret JSON dict")
	rootCmd.MarkFlagRequired("tokenKey")

}
