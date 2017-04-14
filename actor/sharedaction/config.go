package sharedaction

import "code.cloudfoundry.org/cli/util/configv3"

//go:generate counterfeiter . Config

// Config a way of getting basic CF configuration
type Config interface {
	AccessToken() string
	BinaryName() string
	HasTargetedOrganization() bool
	HasTargetedSpace() bool
	PluginHome() string
	Plugins() map[string]configv3.Plugin
	RefreshToken() string
	RemovePlugin(string)
}
