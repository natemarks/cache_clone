package config

// helper functions
import (
	"os"
	"strconv"

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

// GetLogger returns a logger for the application
func GetLogger(verbose bool, s Settings) (log zerolog.Logger) {
	log = zerolog.New(os.Stdout).With().Str("version", version.Version).Timestamp().Logger()
	log = log.Level(zerolog.InfoLevel)
	//log = log.With().Str("SecretID", s.SecretID).Logger()
	//log = log.With().Str("mirror", s.Mirror).Logger()
	//log = log.With().Str("local", s.Local).Logger()
	//log = log.With().Str("remote", s.Remote).Logger()
	if verbose {
		log = log.Level(zerolog.DebugLevel)
		//log = log.With().Str("userKey", s.UserKey).Logger()
		//log = log.With().Str("tokenKey", s.TokenKey).Logger()
		log.Debug().Msg("debug logging enabled")
	}
	return log
}
