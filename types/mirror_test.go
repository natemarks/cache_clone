package types

import (
	"path/filepath"
	"testing"

	"github.com/natemarks/cache_clone/config"
)

// TestClone tests the clone function
func TestClone(t *testing.T) {
	// minimal settings
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
