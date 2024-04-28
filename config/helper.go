package config

// helper functions
import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/natemarks/cache_clone/version"
	"github.com/rs/zerolog"
)

type Settings struct {
	Verbose  bool
	SecretID string
	UserKey  string
	TokenKey string
	Mirror   string
	Local    string
	Remote   string
}

func (s Settings) String() string {
	return "Settings{Verbose: " + strconv.FormatBool(s.Verbose) + ", SecretID: " + s.SecretID + ", UserKey: " + s.UserKey + ", TokenKey: " + s.TokenKey + ", Mirror: " + s.Mirror + ", Local: " + s.Local + ", Remote: " + s.Remote + "}"
}

// MirrorPath returns the local mirror path
// mirror root + remote host and path (without protocol) + repo name
// /home/nmarks/tmp/deleteme.j65Rr2/mirror + stash.imprivata.com/scm/cor_ng + ng.git
func (s Settings) MirrorPath() string {
	remoteHost, remotePath, err := URLHostAndPath(s.Remote)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	return filepath.Join(s.Mirror, remoteHost, remotePath)
}

// GetLogger returns a logger for the application
func GetLogger(s Settings) (log zerolog.Logger) {
	log = zerolog.New(os.Stdout).With().Str("version", version.Version).Timestamp().Logger()
	log = log.Level(zerolog.InfoLevel)
	log = log.With().Str("SecretID", s.SecretID).Logger()
	log = log.With().Str("mirror", s.Mirror).Logger()
	log = log.With().Str("local", s.Local).Logger()
	log = log.With().Str("remote", s.Remote).Logger()
	if s.Verbose {
		log = log.Level(zerolog.DebugLevel)
		log = log.With().Str("userKey", s.UserKey).Logger()
		log = log.With().Str("tokenKey", s.TokenKey).Logger()
		log.Debug().Msg("debug logging enabled")
	}
	return log
}

// JoinPaths joins the elements into a path
func JoinPaths(elements ...string) string {
	return filepath.Join(elements...)
}

// TouchFile creates a file at the path
func TouchFile(path string) error {
	_, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return err
}

// Sha256sum return sh256sum of a string
func Sha256sum(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
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

// Run Runs a shell command and waits to return the results
func Run(c []string) (result Result, err error) {
	var args []string
	baseCommand := c[0]
	args = append(args, c[1:]...)
	cmd := exec.Command(baseCommand, args...)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Error().Err(err).Msg("Error creating stdout pipe")
		return Result{}, err
	}
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Error().Err(err).Msg("Error creating stderr pipe")
		return Result{}, err
	}
	err = cmd.Start()
	if err != nil {
		log.Error().Err(err).Msg("Error starting command")
		return Result{}, err
	}
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
	return Result{cmd.ProcessState.ExitCode(), stdout, stderr}, err

}

// URLHostAndPath given a url return the host and path
//
//	Throw an error if either required value is empty
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
