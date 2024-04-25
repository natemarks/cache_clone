// Package utility provides build and testing tools
package utility

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"time"
)

// DirExists return true if the directory exists
func DirExists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}

// URLHostAndPath given a url return the host and path
//  Throw an error if either required value is empty
func URLHostAndPath(urlString string) (host string, path string, err error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", "", err
	}
	if u.Host == "" || u.Path == "" {
		return "", "", fmt.Errorf("unable to determine host or path from: %s", urlString)
	}
	return u.Host, u.Path, nil
}

// MakeParentDir given a repo location, make sure the parent dirs exist
// so the clone doesn't fail
// example: given /my/clone/dir, run mkdir -p /my/clone
func MakeParentDir(p string) (parentDir string, err error) {
	parent := path.Dir(p)

	result, err := Run([]string{"mkdir", "-p", parent})
	if err != nil {
		return parent, err
	}
	if result.ReturnCode != 0 {
		return parent, errors.New(result.StdErr)
	}
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		return parent, fmt.Errorf("parent directory doesn't exist after mkdir: %s", parent)
	}

	return parent, err
}

// Result is the return from a shell command
type Result struct {
	ReturnCode int
	StdOut     string
	StdErr     string
}

func (r Result) String() string {
	return fmt.Sprintf("Return Code: %d StdOut: %s StdErr: %s", r.ReturnCode, r.StdOut, r.StdErr)
}

func checkErr(e error) {
	if e != nil {
		log.Fatal().Err(e).Msg(e.Error())
	}

}

// Run Runs a shell command and waits to return the results
func Run(c []string) (result Result, err error) {
	var args []string
	baseCommand := c[0]
	args = append(args, c[1:]...)
	cmd := exec.Command(baseCommand, args...)
	outPipe, err := cmd.StdoutPipe()
	checkErr(err)
	errPipe, err := cmd.StderrPipe()
	checkErr(err)
	err = cmd.Start()
	checkErr(err)

	oBuf := new(bytes.Buffer)
	_, err = oBuf.ReadFrom(outPipe)
	if err != nil {
		return Result{}, err
	}
	stdout := oBuf.String()

	eBuf := new(bytes.Buffer)
	_, err = eBuf.ReadFrom(errPipe)
	if err != nil {
		return Result{}, err
	}
	stderr := eBuf.String()

	err = cmd.Wait()
	checkErr(err)

	return Result{cmd.ProcessState.ExitCode(), stdout, stderr}, err

}

// GetTime get time string in preferred format
func GetTime() string {
	currentTime := time.Now()
	return currentTime.Format("20060102-150405")
}

// UpdateRepo Create and checkout a branch with a change so git push has something to do
// Used for testing
func UpdateRepo(repoPath string, testData string) error {
	// Checkout a new branch for the current test run
	result, err := Run([]string{"git", "-C", repoPath, "checkout", "-b", testData})
	if result.ReturnCode != 0 {
		log.Fatal().Err(err).Msgf("Unable to checkout new branch in repo: %s", repoPath)
	}

	val := testData
	data := []byte(val)

	err = ioutil.WriteFile(path.Join(repoPath, testData), data, 0777)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create file: %s", path.Join(repoPath, testData))
	}
	// Add the new file to the repo
	result, _ = Run([]string{"git", "-C", repoPath, "add", "-A"})
	if result.ReturnCode != 0 {
		log.Fatal().Err(err).Msgf("Failed to add file(%s) to repo: %s", path.Join(repoPath, testData), repoPath)
	}
	// Commit the repo change
	result, _ = Run([]string{"git", "-C", repoPath, "commit", "-am", fmt.Sprintf("%s", testData)})
	if result.ReturnCode != 0 {
		log.Fatal().Err(err).Msgf("Failed commit repo: %s", repoPath)
	}
	return nil
}
