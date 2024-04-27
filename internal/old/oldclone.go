package old

import (
	"github.com/natemarks/cache_clone/internal/aws"
	"github.com/natemarks/cache_clone/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	cmd.rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "CLone a remote repo to a local directory using a local mirror",
	Long: `Access the stash credentials from AWS Secret Manager. 
                     Create or update a local mirror of the repo.
                     Clone using the local mirror`,
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

		logger.Info().Msg("Try to Clone Repo")
		err = git.GetMirror(remote, mirror, local, creds.Username, creds.Token, &logger)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Info().Msg("Successfully cloned the repo")
	},
}
