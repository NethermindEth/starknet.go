package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

type Workspace struct {
	Name string `json:"name"`
}

var workspaceCommand = cli.Command{
	Name:    "workspace",
	Aliases: []string{"w"},
	Usage:   "go-starknet workspace",
	Subcommands: []*cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "go-starknet workspace list",
			Action:  workspaceListAction,
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "go-starknet workspace create <workspace>",
			Action:  workspaceCreateAction,
			After:   saveConfiguration,
		},
		{
			Name:    "select",
			Aliases: []string{"s"},
			Usage:   "go-starknet workspace select <workspace>",
			Action:  workspaceSelectAction,
			After:   saveConfiguration,
		},
	},
}

func workspaceListAction(cCtx *cli.Context) error {
	if len(configuration.Workspaces) == 0 {
		configuration.Workspaces = []*Workspace{
			{
				Name: "default",
			},
		}
	}
	if configuration.SelectedWorkspace >= len(configuration.Workspaces) {
		configuration.SelectedWorkspace = 0
		configuration.save()
	}
	for k, workspace := range configuration.Workspaces {
		selected := " "
		if k == configuration.SelectedWorkspace {
			selected = "*"
		}
		fmt.Printf("%s %s\n", selected, workspace.Name)
	}
	return nil
}

func workspaceSelectAction(cCtx *cli.Context) error {
	toSelect := strings.ToLower(cCtx.Args().First())
	if len(configuration.Workspaces) == 0 {
		configuration.Workspaces = []*Workspace{
			{
				Name: "default",
			},
		}
	}
	for k, workspace := range configuration.Workspaces {
		if strings.ToLower(workspace.Name) == toSelect {
			configuration.SelectedWorkspace = k
			return nil
		}
	}
	fmt.Println("workspace not found, use `go-starknet workspace create` first...")
	return errors.New("workspace not found")
}

func workspaceCreateAction(cCtx *cli.Context) error {
	toCreate := strings.ToLower(cCtx.Args().First())
	if len(configuration.Workspaces) == 0 {
		configuration.Workspaces = []*Workspace{
			{
				Name: "default",
			},
		}
	}
	for _, workspace := range configuration.Workspaces {
		if strings.ToLower(workspace.Name) == toCreate {
			return fmt.Errorf("worspace %s already exists", toCreate)
		}
	}
	configuration.Workspaces = append(configuration.Workspaces, &Workspace{Name: toCreate})
	return nil
}
