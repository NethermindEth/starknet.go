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

const profileDirectory = ".go-starknet"
const profileFilename = "profile.json"

type Profile struct {
	DefaultFormat string `json:"defaultFormat,omitempty"`
}

var profileCommand = cli.Command{
	Name:    "profile",
	Aliases: []string{"p"},
	Usage:   "manage the user profile",
	Subcommands: []*cli.Command{
		{
			Name:   "list",
			Usage:  "go-starknet profile list",
			Action: profileList,
		},
		{
			Name:   "set",
			Usage:  "go-starknet profile set name=value",
			Action: profileSet,
		},
	},
}

func initOrLoadProfile() (*Profile, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	profileFullDirectory := filepath.Join(home, profileDirectory)
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
	profileFullFilename := filepath.Join(profileFullDirectory, profileFilename)
	content, err := os.ReadFile(profileFullFilename)
	p := Profile{}
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

func (p Profile) save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	profileFullFilename := filepath.Join(home, profileDirectory, profileFilename)
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

func profileList(cCtx *cli.Context) error {
	p, err := initOrLoadProfile()
	if err != nil {
		return err
	}
	fmt.Printf("profile\n")
	fmt.Printf("  format:   %s\n", or(p.DefaultFormat, "friendly"))
	return nil
}

func profileSet(cCtx *cli.Context) error {
	p, err := initOrLoadProfile()
	if err != nil {
		return err
	}
	values := cCtx.Args()
	if len(values.Slice()) != 1 {
		fmt.Printf("define a value %+v\n", values)
		os.Exit(1)
	}
	changed := false
	k, v, ok := strings.Cut(values.First(), "=")
	if ok && strings.ToLower(k) == "format" {
		p.DefaultFormat = v
		changed = true
	}
	if changed {
		p.save()
	}
	return nil
}
