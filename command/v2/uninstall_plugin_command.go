package v2

import (
	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/v2/shared"
)

type UninstallPluginCommand struct {
	RequiredArgs    flag.PluginName `positional-args:"yes"`
	usage           interface{}     `usage:"CF_NAME uninstall-plugin PLUGIN-NAME"`
	relatedCommands interface{}     `related_commands:"plugins"`

	Config      command.Config
	UI          command.UI
	SharedActor command.SharedActor
}

func (cmd *UninstallPluginCommand) Setup(config command.Config, ui command.UI) error {
	cmd.Config = config
	cmd.UI = ui

	cmd.SharedActor = sharedaction.NewActor()

	return nil
}

func (cmd UninstallPluginCommand) Execute(args []string) error {
	cmd.UI.DisplayTextWithFlavor("Uninstalling plugin {{.PluginName}}...",
		map[string]interface{}{
			"PluginName": cmd.RequiredArgs.PluginName,
		})

	pluginUninstaller := shared.NewPluginUninstaller(cmd.Config, cmd.UI)
	err := cmd.SharedActor.UninstallPlugin(cmd.Config, pluginUninstaller, cmd.RequiredArgs.PluginName)

	if err != nil {
		return shared.HandleError(err)
	}

	cmd.UI.DisplayText("Plugin {{.PluginName}} successfully uninstalled.",
		map[string]interface{}{
			"PluginName": cmd.RequiredArgs.PluginName,
		})

	return nil
}
