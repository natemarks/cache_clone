package git

import (
	"github.com/natemarks/cache_clone/internal/aws"
	"github.com/natemarks/cache_clone/internal/utility"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func Test_getAndPushMirror(t *testing.T) {
	dir, _ := ioutil.TempDir("", "clone_cache_test")
	logger := log.With().Str("test_key", "test_value").Logger()

	creds, err := aws.GetRemoteCredentials(aws.GetRemoteCredentialsInput{
		AWSSMSecretID: os.Getenv("AWSSMSECRETID"),
		UsernameKey:   os.Getenv("USERNAMEKEY"),
		TokenKey:      os.Getenv("TOKENKEY"),
	}, &logger)
	if err != nil {
		t.Error(err)
		return
	}
	type args struct {
		remote         string
		mirror         string
		localParent    string
		remoteUsername string
		remoteToken    string
		logger         *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid_clone", args: args{
			remote:         os.Getenv("REMOTE"),
			mirror:         path.Join(dir, "cache_clone_mirror"),
			localParent:    path.Join(dir, "cache_clone_local"),
			remoteUsername: creds.Username,
			remoteToken:    creds.Token,
			logger:         &logger,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetMirror(tt.args.remote, tt.args.mirror, tt.args.localParent, tt.args.remoteUsername, tt.args.remoteToken, tt.args.logger); (err != nil) != tt.wantErr {
				t.Errorf("GetMirror() error = %v, wantErr %v", err, tt.wantErr)
			}
			err := utility.UpdateRepo(tt.args.localParent, utility.GetTime())
			if err != nil {
				t.Error("Unable to write a change to the repo")
			}
			if err := PushMirror(tt.args.remote, tt.args.mirror, tt.args.localParent, tt.args.remoteUsername, tt.args.remoteToken, tt.args.logger); (err != nil) != tt.wantErr {
				t.Errorf("GetMirror() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetMirror(t *testing.T) {
	dir, _ := ioutil.TempDir("", "clone_cache_test")
	logger := log.With().Str("test_key", "test_value").Logger()

	creds, err := aws.GetRemoteCredentials(aws.GetRemoteCredentialsInput{
		AWSSMSecretID: os.Getenv("AWSSMSECRETID"),
		UsernameKey:   os.Getenv("USERNAMEKEY"),
		TokenKey:      os.Getenv("TOKENKEY"),
	}, &logger)
	if err != nil {
		t.Error(err)
	}

	type args struct {
		remote         string
		mirror         string
		local          string
		remoteUsername string
		remoteToken    string
		log            *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{name: "valid_clone", args: args{
			remote:         os.Getenv("REMOTE"),
			mirror:         path.Join(dir, "cache_clone_mirror"),
			local:          path.Join(dir, "cache_clone_local"),
			remoteUsername: creds.Username,
			remoteToken:    creds.Token,
			log:            &logger,
		}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetMirror(tt.args.remote, tt.args.mirror, tt.args.local, tt.args.remoteUsername, tt.args.remoteToken, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("GetMirror() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
