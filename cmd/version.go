package cmd

import (
	"fmt"

	ver "github.com/natemarks/cache_clone/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cache_clone",
	Long:  `All software has versions. This is cache_clone's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version: ", ver.Version)
	},
}
