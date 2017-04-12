package v2_test

import (
	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command/commandfakes"
	. "code.cloudfoundry.org/cli/command/v2"
	"code.cloudfoundry.org/cli/command/v2/shared"
	"code.cloudfoundry.org/cli/command/v2/v2fakes"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("uninstall-plugin command", func() {
	var (
		cmd        UninstallPluginCommand
		testUI     *ui.UI
		fakeConfig *commandfakes.FakeConfig
		fakeActor  *v2fakes.FakeUninstallPluginActor
		executeErr error
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeActor = new(v2fakes.FakeUninstallPluginActor)

		cmd = UninstallPluginCommand{
			UI:     testUI,
			Config: fakeConfig,
			Actor:  fakeActor,
		}

		cmd.RequiredArgs.PluginName = "some-plugin"
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	Context("when the plugin is installed", func() {
	})

	Context("when the plugin is not installed", func() {
		BeforeEach(func() {
			fakeActor.UninstallPluginReturns(
				v2action.Warnings{"warning-1", "warning-2"},
				v2action.PluginNotFoundError{
					Name: "some-plugin",
				},
			)
		})

		It("outputs warnings and returns an error", func() {
			Expect(testUI.Out).To(Say("Uninstalling plugin some-plugin..."))
			Expect(testUI.Err).To(Say("warning-1"))
			Expect(testUI.Err).To(Say("warning-2"))
			Expect(executeErr).To(MatchError(shared.PluginNotFoundError{
				Name: "some-plugin",
			}))
		})
	})

	Context("when the actor returns an error", func() {
	})
})
