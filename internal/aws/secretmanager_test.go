package aws

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// IsValid checks the GetRemoteCredentialsOutput and returns true if it looks ok
func IsValid(output GetRemoteCredentialsOutput) bool {
	if len(output.SecretSha256sum) != 64 {
		return false
	}
	if len(output.TokenSha256sum) != 64 {
		return false
	}
	if len(output.UsernameSha256sum) != 64 {
		return false
	}
	if output.Token == "" {
		return false
	}
	if output.Username == "" {
		return false
	}
	return true
}

// TestGetRemoteCredentials This test gets credentials based on  a secret manager config proided through env vars
// It executes a loose check on the returned output
func TestGetRemoteCredentials(t *testing.T) {
	t.Skip("skipping test")
	logger := log.With().Str("test_key", "test_value").Logger()
	type args struct {
		i   GetRemoteCredentialsInput
		log *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid", args: args{
			i: GetRemoteCredentialsInput{
				AWSSMSecretID: os.Getenv("AWSSMSECRETID"),
				UsernameKey:   os.Getenv("USERNAMEKEY"),
				TokenKey:      os.Getenv("TOKENKEY"),
			},
			log: &logger,
		},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRemoteCredentials(tt.args.i, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemoteCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !IsValid(got) {
				t.Errorf("GetRemoteCredentials() output doesn't look right")
			}
		})
	}
}
