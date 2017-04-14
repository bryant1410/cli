package command

import "code.cloudfoundry.org/cli/actor/sharedaction"

//go:generate counterfeiter . SharedActor

type SharedActor interface {
	CheckTarget(config sharedaction.Config, targetedOrganizationRequired bool, targetedSpaceRequired bool) error
	UninstallPlugin(config sharedaction.Config, uninstaller sharedaction.PluginUninstaller, name string) error
}
