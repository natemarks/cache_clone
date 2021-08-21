package git

import (
	"errors"
	"fmt"
	"github.com/natemarks/cache_clone/internal/utility"
	"github.com/rs/zerolog"
	"net/url"
	"path"
	"strings"
)

// GetMirror creates a bare mirror from a remote repo, then creates the local working repo form the mirror
// The remote will look something like:
// https://my.git.host/scm/group/project.git
// The mirror input will look like:
// /users/nate/mirror
// The final path to the bare mirror repo would be :
// users/nate/mirror/my.git.host/scm/group/project.git
// The local working repo will be in the path explicitly provided from the command line
// so the user (build job?) knows where to find project files
func GetMirror(remote string, mirror string, local string, remoteUsername string, remoteToken string, log *zerolog.Logger) error {
	// create the mirror path mirror
	remoteHost, remotePath, err := utility.UrlHostAndPath(remote)
	if err != nil {
		log.Fatal().Err(err)
	}
	remoteUrl, err := url.Parse(remote)
	if err != nil {
		return err
	}
	remoteUrl.User = url.UserPassword(remoteUsername, remoteToken)
	mirrorDir := path.Join(mirror, remoteHost, remotePath)
	mirrorParent, err := utility.MakeParentDir(mirrorDir)
	log.Debug().Msgf("Ensuring mirror parent path: %s", mirrorParent)
	_, err = utility.Run([]string{"mkdir", "-p", mirrorParent})
	if err != nil {
		return err
	}
	result, _ := utility.Run([]string{"git", "-C", mirrorDir, "rev-parse", "--is-bare-repository"})
	if result.ReturnCode != 0 || result.StdOut != "true\n" {
		//The local mirror doesn't exist and needs to be cloned
		log.Debug().Msgf("Cloning mirror to : %s", mirrorDir)
		result, err = utility.Run([]string{"git", "-C", mirrorParent, "clone", "--mirror", fmt.Sprint(remoteUrl)})
		if err != nil {
			log.With().Str("errorlevel",
				fmt.Sprint(result.ReturnCode)).Str("stdout", result.StdOut).Str("stderr", result.StdErr)
			log.Fatal().Err(err)
		}
	} else {
		//The local mirror does  exist and needs to be pulled
		log.Debug().Msgf("Mirror exists at : %s. Pulling latest", mirrorDir)
		result, err = utility.Run([]string{"git", "-C", mirrorDir, "fetch", "--all"})
		if err != nil {
			log.With().Str("errorlevel",
				fmt.Sprint(result.ReturnCode)).Str("stdout", result.StdOut).Str("stderr", result.StdErr)
			log.Fatal().Err(err)
		}
	}
	// The mirror is all set , no create the local repo form the local mirror
	// If the local path already exists , error out.
	if utility.DirExists(local) {
		log.Fatal().Err(err).Msgf("Local repo directory already exists: %s", local)
	}
	localParent, err := utility.MakeParentDir(local)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to create local parent: %s", localParent)
	}
	result, err = utility.Run([]string{"git", "clone", mirrorDir, local})
	log.Debug().Msgf("Creating local clone(%s) from mirror(%s)", local, mirrorDir)
	if err != nil {
		log.With().Str("errorlevel",
			fmt.Sprint(result.ReturnCode)).Str("stdout", result.StdOut).Str("stderr", result.StdErr)
		log.Fatal().Err(err)
	}

	return err
}

// PushMirror pushes the local repo to the local mirror, then the mirror to the remote
func PushMirror(remote string, mirror string, local string, remoteUsername string, remoteToken string, log *zerolog.Logger) error {
	// create the mirror path mirror
	remoteHost, remotePath, err := utility.UrlHostAndPath(remote)
	if err != nil {
		log.Fatal().Err(err)
	}
	remoteUrl, err := url.Parse(remote)
	if err != nil {
		return err
	}
	remoteUrl.User = url.UserPassword(remoteUsername, remoteToken)
	mirrorDir := path.Join(mirror, remoteHost, remotePath)
	// Check the status of the local repo before trying to push
	log.Info().Msgf("Checking status of local repo: %s", local)
	result, _ := utility.Run([]string{"git", "-C", local, "status", "--short"})
	if err != nil {
		log.Fatal().Err(err)
	}
	// git status --short stdout will be empty if the repo is clean
	if result.StdOut != "" {
		msg := fmt.Sprintf("Unable to push dirty repo: %s", local)
		err := errors.New(msg)
		log.Fatal().Err(err)
	}
	// Get the current branch name so we can push it
	log.Info().Msgf("Get current branch of local repo: %s", local)
	result, _ = utility.Run([]string{"git", "-C", local, "branch", "--show-current"})
	branch := strings.TrimSuffix(result.StdOut, "\n")
	log.Debug().Msgf("Got current branch of local repo (%s): %s", local, branch)

	//Push the current local branch to the mirror
	log.Info().Msgf("Pushing local repo(%s) to mirror(%s)", local, mirrorDir)
	result, _ = utility.Run([]string{"git", "-C", local, "push", "--set-upstream", "origin", branch})
	if result.ReturnCode != 0 {
		msg := fmt.Sprintf("Unable to local repo (%s) to mirror (%s)", local, mirrorDir)
		err = errors.New(msg)
		log.Fatal().Err(err)
	}
	// Push the mirro to the remote
	log.Info().Msgf("Pushing mirror(%s) to remote(%s)", mirrorDir, remote)
	result, _ = utility.Run([]string{"git", "-C", mirrorDir, "push"})
	if result.ReturnCode != 0 {
		msg := fmt.Sprintf("Unable to push mirror (%s) to remote (%s)", mirrorDir, remote)
		err = errors.New(msg)
		log.Fatal().Err(err)
	}
	return nil
}
