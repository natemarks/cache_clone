package cmd

import (
	"github.com/natemarks/cache_clone/config"
	"github.com/natemarks/cache_clone/types"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push build repo changes through the local mirror to remote",
	Long: `Access the stash credentials from AWS Secret Manager. 
                     Push the build repo changes to the local mirror.
                     Push the local mirror to the remote`,
	Run: func(cmd *cobra.Command, args []string) {
		log := config.GetLogger(settings)
		//TODO: remove the next three lines. I don't think we need cred for push
		//log.Debug().Msg("Getting credentials from AWS Secret Manager")
		//r := types.NewHTTPSRemote(settings.Remote)
		//creds := *types.NewCredential(settings, &log)
		types.PushMirror(settings, &log)

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
