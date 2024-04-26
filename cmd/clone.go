/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/natemarks/cache_clone/config"
	"github.com/natemarks/cache_clone/internal/aws"
	"github.com/natemarks/cache_clone/internal/git"
	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a remote repo to a local directory using a local mirror",
	Long: `Access the stash credentials from AWS Secret Manager. 
                     Create or update a local mirror of the repo.
                     Clone using the local mirror`,
	Run: func(cmd *cobra.Command, args []string) {
		log := config.GetLogger(settings)
		// Get credentials from AWS Secret Manager
		creds, err := aws.GetRemoteCredentials(aws.GetRemoteCredentialsInput{
			AWSSMSecretID: settings.SecretID,
			UsernameKey:   settings.UserKey,
			TokenKey:      settings.TokenKey,
		}, &log)
		log.Info().Msg("Try to Clone Repo")
		err = git.GetMirror(settings.Remote, settings.Mirror, settings.Local, creds.Username, creds.Token, &log)
		if err != nil {
			log.Fatal().Err(err)
		}
		log.Info().Msg("Successfully cloned the repo")
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
