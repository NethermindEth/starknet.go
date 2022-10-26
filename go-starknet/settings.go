package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const configDirectory = ".go-starknet"
const configFilename = "configuration.json"

type Configuration struct {
	DefaultFormat     string       `json:"defaultFormat,omitempty"`
	SelectedWorkspace int          `json:"selectedWorkspace,omitempty"`
	Workspaces        []*Workspace `json:"workspaces,omitempty"`
}

var configuration *Configuration

var settingsCommand = cli.Command{
	Name:    "settings",
	Aliases: []string{"s"},
	Usage:   "go-starknet settings",
	Subcommands: []*cli.Command{
		{
			Name:   "list",
			Usage:  "go-starknet settings list",
			Action: settingsList,
		},
		{
			Name:   "set",
			Usage:  "go-starknet settings set name=value",
			Action: settingsSet,
		},
	},
}

func saveConfiguration(cCtx *cli.Context) error {
	return configuration.save()
}

func initOrLoadConfig() (*Configuration, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	profileFullDirectory := filepath.Join(home, configDirectory)
	v, err := os.Stat(profileFullDirectory)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(profileFullDirectory, 0755)
		if err != nil {
			return nil, err
		}
		v, err = os.Stat(profileFullDirectory)
	}
	if err != nil {
		return nil, err
	}
	if !v.IsDir() {
		return nil, fmt.Errorf("%s not directory", v.Name())
	}
	profileFullFilename := filepath.Join(profileFullDirectory, configFilename)
	content, err := os.ReadFile(profileFullFilename)
	p := Configuration{}
	if err != nil && errors.Is(err, os.ErrNotExist) {
		content, err = json.MarshalIndent(p, " ", "  ")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(profileFullFilename, content, 0755)
		if err != nil {
			return nil, err
		}
		return &p, nil
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &p)
	return &p, err
}

func (p Configuration) save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	profileFullFilename := filepath.Join(home, configDirectory, configFilename)
	content, err := json.MarshalIndent(p, " ", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(profileFullFilename, content, 0755)
}

func or(a string, b string) string {
	if a != "" {
		return a
	}
	return b
}

func settingsList(cCtx *cli.Context) error {
	fmt.Printf("current settings\n")
	fmt.Printf("  format:   %s\n", or(configuration.DefaultFormat, "friendly"))
	return nil
}

func settingsSet(cCtx *cli.Context) error {
	values := cCtx.Args()
	if len(values.Slice()) != 1 {
		fmt.Printf("define a value %+v\n", values)
		os.Exit(1)
	}
	changed := false
	k, v, ok := strings.Cut(values.First(), "=")
	if ok && strings.ToLower(k) == "format" {
		switch strings.ToLower(v) {
		case "friendly":
			configuration.DefaultFormat = "friendly"
			changed = true
		case "raw":
			configuration.DefaultFormat = "raw"
			changed = true
		default:
			fmt.Println("unsupported format, should be friendly or raw")
			os.Exit(1)
		}
	}
	if changed {
		configuration.save()
	}
	return nil
}
