package cmd

import (
	"github.com/natemarks/cache_clone/internal/aws"
	"github.com/natemarks/cache_clone/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push build repo changes through the local mirror to remote",
	Long: `Access the stash credentials from AWS Secret Manager. 
                     Push the build repo changes to the local mirror.
                     Push the local mirro to the remote`,
	Run: func(cmd *cobra.Command, args []string) {
		// set up the logger
		logger := log.With().Str("SecretID", secretID).Logger()
		logger = logger.With().Str("userKey", userKey).Logger()
		logger = logger.With().Str("tokenKey", tokenKey).Logger()
		logger = logger.With().Str("mirror", mirror).Logger()
		logger = logger.With().Str("local", local).Logger()
		logger = logger.With().Str("remote", remote).Logger()
		// Get credentials from AWS Secret Manager
		creds, err := aws.GetRemoteCredentials(aws.GetRemoteCredentialsInput{
			AWSSMSecretID: secretID,
			UsernameKey:   userKey,
			TokenKey:      tokenKey,
		}, &logger)

		logger.Info().Msg("Try to Push Repo")
		err = git.PushMirror(remote, mirror, local, creds.Username, creds.Token, &logger)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Info().Msg("Successfully Pushed Repo")
	},
}
