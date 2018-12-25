package types

type ProtocolDefinition struct {
	Version 	uint64	`json:"version"`
	Software	string	`json:"software"`
	Height		uint64	`json:"height"`
}

func (pd ProtocolDefinition) GetVersion() uint64 {
	return pd.Version
}

func (pd ProtocolDefinition) GetSoftware() string {
	return pd.Software
}

func (pd ProtocolDefinition) GetHeight() uint64 {
	return pd.Height
}