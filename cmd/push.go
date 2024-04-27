/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/natemarks/cache_clone/config"
	"github.com/natemarks/cache_clone/internal/utility"
	"github.com/natemarks/cache_clone/types"
	"strings"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := config.GetLogger(settings)
		log.Debug().Msg("Getting credentials from AWS Secret Manager")
		//r := types.NewHTTPSRemote(settings.Remote)
		//creds := *types.NewCredential(settings, &log)
		mirror := *types.NewMirror(settings, &log)
		log.Info().Msgf("Checking status of local repo: %s", settings.Local)
		result, err := utility.Run([]string{"git", "-C", settings.Local, "status", "--short"})
		if err != nil || result.ReturnCode != 0 || result.StdOut != "" {
			log.Error().Msgf("Unable to push dirty repo: %s", settings.Local)
			log.Fatal().Err(err).Msg(result.String())
		}

		// Get the current branch name so we can push it
		log.Debug().Msgf("Get current branch of local repo: %s", settings.Local)
		result, _ = utility.Run([]string{"git", "-C", settings.Local, "branch", "--show-current"})
		branch := strings.TrimSuffix(result.StdOut, "\n")
		log.Info().Msgf("Got current branch of local repo (%s): %s", settings.Local, branch)

		//Push the current local branch to the mirror
		log.Debug().Msgf("Pushing local repo(%s) to mirror(%s)", settings.Local, mirror.Path)
		result, err = utility.Run([]string{"git", "-C", settings.Local, "push", "--set-upstream", "origin", branch})
		if err != nil || result.ReturnCode != 0 || result.StdOut != "" {
			log.Error().Msgf("Unable to push local repo (%s) to mirror (%s)", settings.Local, mirror.Path)
			log.Fatal().Err(err).Msg(result.String())
		}

		// Push the mirror to the remote
		log.Debug().Msgf("Pushing mirror(%s) to remote(%s)", mirror.Path, settings.Remote)
		result, _ = utility.Run([]string{"git", "-C", mirror.Path, "push"})
		if err != nil || result.ReturnCode != 0 || result.StdOut != "" {
			log.Error().Msgf("Unable to push mirror (%s) to remote (%s)", mirror.Path, settings.Remote)
			log.Fatal().Err(err).Msg(result.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
