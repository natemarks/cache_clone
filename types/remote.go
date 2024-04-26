package types

// TODO: simplify remote functions by getting rid of the struct and just using the functions
import (
	"fmt"
	"net/url"
)

// HTTPSRemote is a struct that represents a remote git repository
type HTTPSRemote struct {
	Host string // https://my.git.host/my/git/repo.git -> my.git.host
	Path string // https://my.git.host/my/git/repo.git -> my/git/repo.git
	URL  *url.URL
}

// ConnectionString returns the connection string for the remote
func (r HTTPSRemote) ConnectionString(c Credential) string {
	r.URL.User = url.UserPassword(c.Username, c.Token)
	return r.URL.String()
}

// NewHTTPSRemote returns a HTTPSRemote struct from a remote URL
func NewHTTPSRemote(remoteURL string) *HTTPSRemote {
	u, err := url.Parse(remoteURL)
	if err != nil {
		panic(err)
	}
	if u.Host == "" || u.Path == "" {
		panic(fmt.Errorf("unable to determine host or path from: %s", remoteURL))
	}
	ff, err := url.Parse(remoteURL)
	if err != nil {
		panic(err)
	}
	return &HTTPSRemote{
		Host: u.Host,
		Path: u.Path,
		URL:  ff,
	}
}
