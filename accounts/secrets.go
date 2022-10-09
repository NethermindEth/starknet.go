package main

import (
	"encoding/json"
	"os"
)

const SECRET_FILE_NAME = ".starknet-account.json"

type accountPlugin struct {
	PrivateKey       string `json:"privateKey"`
	PublicKey        string `json:"publicKey"`
	PluginClassHash  string `json:"pluginClassHash,omitempty"`
	AccountClassHash string `json:"accountClassHash,omitempty"`
	AccountAddress   string `json:"accountAddress"`
	Version          string `json:"accountVersion"`
	Plugin           bool   `json:"accountPlugin"`
}

func (ap *accountPlugin) Read(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	json.Unmarshal(content, ap)
	return nil
}

func (ap *accountPlugin) Write(filename string) error {
	content, err := json.Marshal(ap)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, content, 0664)
}
