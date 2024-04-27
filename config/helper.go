package config

// helper functions
import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/natemarks/cache_clone/internal/utility"
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
	remoteHost, remotePath, err := utility.URLHostAndPath(s.Remote)
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
