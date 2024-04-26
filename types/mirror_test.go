package types

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/natemarks/cache_clone/config"
)

const (
	testFile = "deleteme.txt"
)

// set the time as a string, so it's different every time the test runs
// but  consistent throughout the test
var currentTime = epochTimeString()

// epochTimeString returns the current Unix epoch time in seconds as a string
func epochTimeString() string {
	// Get the current Unix epoch time in seconds
	epochTime := time.Now().Unix()

	// Convert the epoch time to a string
	epochTimeString := strconv.FormatInt(epochTime, 10)

	return epochTimeString
}

// writeStringToFile writes a string to a file
// use this to create a test commit
func writeStringToFile(filePath string) error {
	data := " This file is added as a cache_clone test. feel free to delete it\n"
	data += currentTime
	// Write the string data to the specified file path
	err := os.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

// checkoutNewBranch checks out a test branch
func checkoutNewBranch(s config.Settings, log *zerolog.Logger) {
	result, err := config.Run(
		[]string{"git", "-C", s.Local, "checkout", "-b", fmt.Sprintf("cache_clone_%s", currentTime)})
	if err != nil || result.ReturnCode != 0 || result.StdOut != "" {
		log.Error().Msgf("Unable to checkout new branch: %s", result.String())
		log.Fatal().Err(err).Msg(result.String())
	}
}

// commitNewBranch commits a test branch
func commitNewBranch(s config.Settings, log *zerolog.Logger) {
	result, err := config.Run(
		[]string{"git", "-C", s.Local, "add", testFile})
	if err != nil || result.ReturnCode != 0 || result.StdOut != "" {
		log.Error().Msgf("Unable to add new branch: %s", result.String())
		log.Fatal().Err(err).Msg(result.String())
	}
	result, err = config.Run(
		[]string{"git", "-C", s.Local, "commit", "-m", fmt.Sprintf("cache_clone_%s", currentTime)})
	if err != nil || result.ReturnCode != 0 || result.StdErr != "" {
		log.Error().Msgf("Unable to commit new branch: %s", result.String())
		log.Fatal().Err(err).Msg(result.String())
	}
}

// TestClone tests the clone function
// This test covers the clone function.
func TestClone(t *testing.T) {
	// minimal settings. no credentials needed because we test cloning from a public repo
	s := config.Settings{
		Verbose: true,
		Mirror:  t.TempDir(),
		Remote:  "https://github.com/natemarks/cache_clone.git",
	}
	// get a logger
	log := config.GetLogger(s)

	// create a new mirror
	m := NewMirror(s, &log)
	// confirm the mirror is not cloned
	if m.CheckClone(&log) {
		t.Fatalf("Mirror should not be cloned yet")
	}
	// create the mirror
	// this works with empty credentials because the repo is public
	// and obviates the need for AWS Secret Manager access
	m.CreateClone(*NewHTTPSRemote(s.Remote), Credential{}, &log)
	// confirm the mirror is cloned
	if !m.CheckClone(&log) {
		t.Fatalf("Mirror should be cloned")

	}
	// update the mirror
	m.UpdateClone(&log)
	// clone the mirror locally
	m.MakeLocal(filepath.Join(s.Mirror, "local"), &log)
}

// TestCloneAndPush tests the push function
// it covers clone and push by cloning, then creating a test branch and pushing the test branch
func TestCloneAndPush(t *testing.T) {
	t.Skip("Skipping test that requires AWS Secret Manager access")
	testDir := t.TempDir()
	// minimal settings
	s := config.Settings{
		Verbose:  true,
		Mirror:   filepath.Join(testDir, "mirrorRoot"),
		Remote:   "https://stash.imprivata.com/scm/cldops/dna.git",
		Local:    filepath.Join(testDir, "local"),
		SecretID: "/azure_agent/vpn-connected-self-managed-agents",
		UserKey:  "azure_agent_username",
		TokenKey: "azure_agent_token",
	}
	// get a logger
	log := config.GetLogger(s)

	// create a new mirror
	m := NewMirror(s, &log)
	// confirm the mirror is not cloned
	if m.CheckClone(&log) {
		t.Fatalf("Mirror should not be cloned yet")
	}
	// create the mirror
	// this works with empty credentials because the repo is public
	// and obviates the need for AWS Secret Manager access
	m.CreateClone(*NewHTTPSRemote(s.Remote), *NewCredential(s, &log), &log)
	// confirm the mirror is cloned
	if !m.CheckClone(&log) {
		t.Fatalf("Mirror should be cloned")

	}
	// update the mirror
	m.UpdateClone(&log)
	// clone the mirror locallyz
	m.MakeLocal(s.Local, &log)
	// checkout a branch for the test
	checkoutNewBranch(s, &log)
	writeStringToFile(filepath.Join(s.Local, testFile))
	commitNewBranch(s, &log)
	PushMirror(s, &log)
}
