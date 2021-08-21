// Package cmd ...
package cmd

import (
	"fmt"
	ver "github.com/natemarks/cache_clone/version"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var log = zerolog.New(os.Stdout).With().Str("version", ver.Version).Timestamp().Logger()
var cfgFile string

// verbose mode
var verbose bool

// remote is the git repo url. example: https://my.git.host/scm/group/project.git
// local the local working directory, usually placed inside a temporary build directory. example:${BUILD_DIR}/project
// mirror is the persistent mirror path on the agent host. Each repo s gets a unique path inside the mirror path based
// on the repo url. example: the the mirror path is ${HOME}/cache_clone_mirror  and the remote is:
// https://my.git.host/scm/group/project.git
// The bare git mirror repo will be: :
//  ${HOME}/cache_clone_mirror/my.git.host/scm/group/project.git
// secretID is the AWS secret manager secret ID path. it's the location of the secret JSON doc
// userKey is the key in the JSON doc associated with the username value for the remote
// tokenKey is the key in the JSON doc associated with the token value for the remote
var mirror, local, remote, secretID, userKey, tokenKey string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cache_clone",
	Short: "Clone a repo using a local mirror",
	Long: `Using credentials stored in AWS Secret Manager, create/update a local mirror.
Then use the local mirror to create local clone`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cache_clone.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&mirror, "mirror", "", "Location for all mirror repos")
	err := rootCmd.MarkPersistentFlagRequired("mirror")
	if err != nil {
		log.Fatal().Err(err)
	}
	rootCmd.PersistentFlags().StringVar(&local, "local", "", "Location to create the repo clone")
	err = rootCmd.MarkPersistentFlagRequired("local")
	if err != nil {
		log.Fatal().Err(err)
	}
	rootCmd.PersistentFlags().StringVar(&remote, "remote", "", "git remote server url. example: https://my.git.com/my/project.git")
	err = rootCmd.MarkPersistentFlagRequired("remote")
	if err != nil {
		log.Fatal().Err(err)
	}
	rootCmd.PersistentFlags().StringVar(&secretID, "secretID", "", "AWS Secret Manager secretID path")
	err = rootCmd.MarkPersistentFlagRequired("secretID")
	if err != nil {
		log.Fatal().Err(err)
	}
	rootCmd.PersistentFlags().StringVar(&userKey, "userKey", "", "username key in the secret JSON dict")
	err = rootCmd.MarkPersistentFlagRequired("userKey")
	if err != nil {
		log.Fatal().Err(err)
	}
	rootCmd.PersistentFlags().StringVar(&tokenKey, "tokenKey", "", "token key in the secret JSON dict")
	err = rootCmd.MarkPersistentFlagRequired("tokenKey")
	if err != nil {
		log.Fatal().Err(err)
	}
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		log.Info().Msg("Starting")
	}
	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {

		log.Info().Msg("Stopping")
	}
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cache_clone" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cache_clone")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
