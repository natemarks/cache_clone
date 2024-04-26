package config

// CreateConfig represents the configuration for the create command
type CreateConfig struct {
	Region    string `json:"region"`
	AccountID string `json:"accountID"`
	Debug     bool   `json:"debug"`
	VPCID     string `json:"vpcID"`
}

// GetCreateConfig returns the configuration for the show command
func GetCreateConfig(vpcID string, debug bool) (config CreateConfig, err error) {
	awsInfo, err := GetAWSInfo()
	if err != nil {
		return config, err

	}
	config = CreateConfig{
		Region:    awsInfo.Region,
		AccountID: awsInfo.AccountID,
		Debug:     debug,
		VPCID:     vpcID,
	}
	return config, nil
}
