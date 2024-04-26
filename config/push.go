package config

// DeleteConfig represents the configuration for the delete command
type DeleteConfig struct {
	Region    string `json:"region"`
	AccountID string `json:"accountID"`
	Debug     bool   `json:"debug"`
	FlowLogID string `json:"flowLogID"`
}

// GetDeleteConfig returns the configuration for the show command
func GetDeleteConfig(flowLogID string, debug bool) (config DeleteConfig, err error) {
	awsInfo, err := GetAWSInfo()
	if err != nil {
		return config, err

	}
	config = DeleteConfig{
		Region:    awsInfo.Region,
		AccountID: awsInfo.AccountID,
		Debug:     debug,
		FlowLogID: flowLogID,
	}
	return config, nil
}
