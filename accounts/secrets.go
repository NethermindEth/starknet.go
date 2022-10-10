package accounts

import (
	"encoding/json"
	"os"
)

type AccountPlugin struct {
	PrivateKey       string `json:"privateKey"`
	PublicKey        string `json:"publicKey"`
	PluginClassHash  string `json:"pluginClassHash,omitempty"`
	AccountClassHash string `json:"accountClassHash,omitempty"`
	AccountAddress   string `json:"accountAddress"`
	Version          string `json:"accountVersion"`
	Plugin           bool   `json:"accountPlugin"`
	filename         string
}

func (ap *AccountPlugin) Read(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	json.Unmarshal(content, ap)
	return nil
}

func (ap *AccountPlugin) Write(filename string) error {
	content, err := json.Marshal(ap)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, content, 0664)
}
