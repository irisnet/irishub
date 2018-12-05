package common

type ProtocolDefinition struct {
	version 	uint	`json:"version"`
	software	string	`json:"software"`
	height		uint	`json:"height"`
}

func (pd ProtocolDefinition) GetVersion() uint {
	return pd.version
}

func (pd ProtocolDefinition) GetSoftware() string {
	return pd.software
}

func (pd ProtocolDefinition) GetHeight() uint {
	return pd.height
}