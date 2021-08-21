package aws

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog"
)

type GetRemoteCredentialsInput struct {
	// The AWS SecretManager secret identifier. ex: /path/to/y/secret
	AWSSMSecretID string
	// Username JSON Key for the secret value to get. ex: my_repo_username
	UsernameKey string
	// Token JSON Key for the secret value to get. ex: my_repo_token
	TokenKey string
}

// GetRemoteCredentialsOutput The output includes the credentiols and the sha256
// some of the username, toke and the AWS SM secret document. the sha256sum makes
// it easier to test and troubleshoot bad values without compromising the secure
// data
type GetRemoteCredentialsOutput struct {
	//The sha256sum of the AWS SM secret json document
	SecretSha256sum string
	// Remote username
	Username string
	// Sha256sum of the remote username
	UsernameSha256sum string
	// Remote token
	Token string
	// Sha256sum of the remote token
	TokenSha256sum string
}

func Sha256sum(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

// GetRemoteCredentials returns the remote credentials and sha256sums
func GetRemoteCredentials(i GetRemoteCredentialsInput, log *zerolog.Logger) (GetRemoteCredentialsOutput, error) {

	// Setup the client
	log.Info().Msg("setting up the AWS Secret Manager client")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err)
	}

	SecretClient := *secretsmanager.NewFromConfig(cfg)

	SecretInput := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(i.AWSSMSecretID),
		VersionId:    nil,
		VersionStage: nil,
	}
	// Get the secret doc from AWS

	log.Info().Msg("getting the secret doc from AWS SM")
	secretDoc, err := SecretClient.GetSecretValue(context.TODO(), SecretInput)
	if err != nil {
		log.Fatal().Err(err)
	}

	// unmarshal the JSON secret doc into a map. If the structure isn't a map this will fail
	log.Info().Msg("Unmarshalling credentials from AWSSM secret doc")
	var objmap map[string]string
	err = json.Unmarshal([]byte(*secretDoc.SecretString), &objmap)
	if err != nil {
		log.Fatal().Err(err)
	}
	// Use the provided username and token key names to get the credential values
	username := objmap[i.UsernameKey]
	token := objmap[i.TokenKey]
	result := GetRemoteCredentialsOutput{
		SecretSha256sum:   Sha256sum(*secretDoc.SecretString),
		Username:          username,
		UsernameSha256sum: Sha256sum(username),
		Token:             token,
		TokenSha256sum:    Sha256sum(token),
	}
	log.Debug().Msgf("SecretJSON Document(sha256): %s", Sha256sum(result.SecretSha256sum))
	log.Debug().Msgf("Username(sha256): %s", Sha256sum(result.Username))
	log.Debug().Msgf("Token(sha256): %s", Sha256sum(result.Token))
	return result, err
}
