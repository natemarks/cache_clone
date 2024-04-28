package types

import (
	"context"
	"encoding/json"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/natemarks/cache_clone/config"
	"github.com/rs/zerolog"
)

// Credential is a struct that represents a credential
// store the sha256sums for logging/debugging purposes
type Credential struct {
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

// NewCredential creates a new Credential struct
func NewCredential(s config.Settings, log *zerolog.Logger) *Credential {

	// Set up the client
	log.Debug().Msg("setting up the AWS Secret Manager client")
	cfg, err := awscfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	SecretClient := *secretsmanager.NewFromConfig(cfg)

	SecretInput := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(s.SecretID),
		VersionId:    nil,
		VersionStage: nil,
	}
	// Get the secret doc from AWS

	log.Debug().Msg("getting the secret doc from AWS SM")
	secretDoc, err := SecretClient.GetSecretValue(context.TODO(), SecretInput)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	// unmarshal the JSON secret doc into a map. If the structure isn't a map this will fail
	log.Debug().Msg("unmarshalling credentials from AWSSM secret doc")
	var objmap map[string]string
	err = json.Unmarshal([]byte(*secretDoc.SecretString), &objmap)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	// Use the provided username and token key names to get the credential values
	doc := *secretDoc.SecretString
	username := objmap[s.UserKey]
	token := objmap[s.TokenKey]

	log.Debug().Msgf("SecretJSON Document(sha256): %s", config.Sha256sum(doc))
	log.Debug().Msgf("Username(sha256): %s", config.Sha256sum(username))
	log.Debug().Msgf("Token(sha256): %s", config.Sha256sum(token))

	return &Credential{
		SecretSha256sum:   config.Sha256sum(doc),
		Username:          username,
		UsernameSha256sum: config.Sha256sum(username),
		Token:             token,
		TokenSha256sum:    config.Sha256sum(token),
	}
}
