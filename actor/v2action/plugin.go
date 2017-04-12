package v2action

import "fmt"

// PluginNotFoundError is an error returned when a plugin is not found.
type PluginNotFoundError struct {
	Name string
}

// Error method to display the error message.
func (e PluginNotFoundError) Error() string {
	return fmt.Sprintf("Plugin name %s does not exist", e.Name)
}
