package v2

import (
	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/v2/shared"
)

//go:generate counterfeiter . UninstallPluginActor
type UninstallPluginActor interface {
	UninstallPlugin(name string) (v2action.Warnings, error)
}

type UninstallPluginCommand struct {
	RequiredArgs    flag.PluginName `positional-args:"yes"`
	usage           interface{}     `usage:"CF_NAME uninstall-plugin PLUGIN-NAME"`
	relatedCommands interface{}     `related_commands:"plugins"`

	Config command.Config
	UI     command.UI
	Actor  UninstallPluginActor
}

func (cmd *UninstallPluginCommand) Setup(config command.Config, ui command.UI) error {
	cmd.Config = config
	cmd.UI = ui

	return nil
}

func (cmd UninstallPluginCommand) Execute(args []string) error {
	cmd.UI.DisplayTextWithFlavor("Uninstalling plugin {{.PluginName}}...",
		map[string]interface{}{
			"PluginName": cmd.RequiredArgs.PluginName,
		})

	warnings, err := cmd.Actor.UninstallPlugin(cmd.RequiredArgs.PluginName)
	cmd.UI.DisplayWarnings(warnings)

	return shared.HandleError(err)
}
