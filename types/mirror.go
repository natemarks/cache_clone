package types

import (
	"path"

	"github.com/natemarks/cache_clone/config"
	"github.com/rs/zerolog"
)

// TODO: simplify mirriring by gettting rid of the struct and just using the functions

// Mirror is a struct that represents a git mirror
// mirrors should never be pulled or cloned more than once
// but tracking it makes it safe to run those functions multiple times
type Mirror struct {
	IsCloned bool
	IsPulled bool
	Path     string
}

// CheckClone returns true if the mirror is cloned
// it also sets the IsCloned flag. Use this to avoid rerunning git commands
func (m Mirror) CheckClone(log *zerolog.Logger) bool {
	// if this is set to true, we don't need to check again
	if m.IsCloned {
		log.Debug().Msgf("already confirmed the mirror is cloned: %s", m.Path)
		return true
	}
	result, _ := config.Run([]string{"git", "-C", m.Path, "rev-parse", "--is-bare-repository"})
	if result.ReturnCode != 0 || result.StdOut != "true\n" {
		log.Debug().Msgf("mirror is not cloned: %s", result.String())
		m.IsCloned = false
		return false
	}
	log.Debug().Msgf("mirror is cloned: %s", m.Path)
	m.IsCloned = true
	return true
}

// CreateClone creates a mirror of a remote repo
func (m Mirror) CreateClone(r HTTPSRemote, c Credential, log *zerolog.Logger) {
	mirrorParent := path.Dir(m.Path)

	result, err := config.Run([]string{"mkdir", "-p", mirrorParent})
	if err != nil || result.ReturnCode != 0 {
		log.Fatal().Err(err).Msg(result.String())
	}
	log.Debug().Msgf("cloning mirror to : %s", m.Path)
	result, err = config.Run([]string{"git", "-C", mirrorParent, "clone", "--mirror", r.ConnectionString(c)})
	if err != nil {
		log.Fatal().Err(err).Msg(result.String())
	}

}

// UpdateClone updates the mirror with the latest changes
func (m Mirror) UpdateClone(log *zerolog.Logger) {
	if m.IsPulled {
		log.Debug().Msgf("mirror is already pulled: %s", m.Path)
		return
	}
	log.Debug().Msgf("mirror exists at : %s. Pulling latest", m.Path)
	result, err := config.Run([]string{"git", "-C", m.Path, "fetch", "--all"})
	if err != nil {
		log.Fatal().Err(err).Msg(result.String())
	}
	m.IsPulled = true
}

// MakeLocal creates a local clone from the mirror
func (m Mirror) MakeLocal(l string, log *zerolog.Logger) {
	localParent := path.Dir(l)
	log.Debug().Msgf("Ensuring local parent path: %s", localParent)
	result, err := config.Run([]string{"mkdir", "-p", localParent})
	if err != nil || result.ReturnCode != 0 {
		log.Fatal().Err(err).Msg(result.String())
	}
	log.Debug().Msgf("Creating local clone(%s) from mirror(%s)", l, m.Path)
	result, err = config.Run([]string{"git", "clone", m.Path, l})
	if err != nil {
		log.Fatal().Err(err).Msg(result.String())
	}
}

// NewMirror returns a new Mirror struct
func NewMirror(s config.Settings, log *zerolog.Logger) *Mirror {
	remote := NewHTTPSRemote(s.Remote)

	return &Mirror{
		IsCloned: false,
		IsPulled: false,
		Path:     config.JoinPaths(s.Mirror, remote.Host, remote.Path),
	}
}
