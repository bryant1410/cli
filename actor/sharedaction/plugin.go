package sharedaction

import (
	"fmt"
	"os"
)

// PluginNotFoundError is an error returned when a plugin is not found.
type PluginNotFoundError struct {
	Name string
}

// Error method to display the error message.
func (e PluginNotFoundError) Error() string {
	return fmt.Sprintf("Plugin name %s does not exist", e.Name)
}

//go:generate counterfeiter . PluginUninstaller
type PluginUninstaller interface {
	Uninstall(location string) error
}

func (_ Actor) UninstallPlugin(config Config, uninstaller PluginUninstaller, name string) error {
	plugin, exist := config.Plugins()[name]
	if !exist {
		return PluginNotFoundError{Name: name}
	}

	err := uninstaller.Uninstall(plugin.Location)
	if err != nil {
		return err
	}

	// sleep for 500 ms???

	err = os.Remove(plugin.Location)
	if err != nil {
		return err
	}

	config.RemovePlugin(name)

	return nil
}
