package v2_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/command/commandfakes"
	. "code.cloudfoundry.org/cli/command/v2"
	"code.cloudfoundry.org/cli/command/v2/shared"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("uninstall-plugin command", func() {
	var (
		cmd             UninstallPluginCommand
		testUI          *ui.UI
		fakeConfig      *commandfakes.FakeConfig
		fakeSharedActor *commandfakes.FakeSharedActor
		executeErr      error
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeSharedActor = new(commandfakes.FakeSharedActor)

		cmd = UninstallPluginCommand{
			UI:          testUI,
			Config:      fakeConfig,
			SharedActor: fakeSharedActor,
		}

		cmd.RequiredArgs.PluginName = "some-plugin"
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	Context("when the plugin is installed", func() {
		BeforeEach(func() {
			fakeSharedActor.UninstallPluginReturns(nil)
		})

		It("outputs warnings and a success message", func() {
			Expect(executeErr).ToNot(HaveOccurred())
			Expect(testUI.Out).To(Say("Uninstalling plugin some-plugin..."))

			Expect(testUI.Out).To(Say("Plugin some-plugin successfully uninstalled."))

			Expect(fakeSharedActor.UninstallPluginCallCount()).To(Equal(1))
			config, _, pluginName := fakeSharedActor.UninstallPluginArgsForCall(0)
			Expect(config).To(Equal(fakeConfig))
			Expect(pluginName).To(Equal("some-plugin"))
		})
	})

	Context("Errors", func() {
		Context("when the plugin is not installed", func() {
			BeforeEach(func() {
				fakeSharedActor.UninstallPluginReturns(
					sharedaction.PluginNotFoundError{
						Name: "some-plugin",
					},
				)
			})

			It("returns an error", func() {
				Expect(testUI.Out).To(Say("Uninstalling plugin some-plugin..."))
				Expect(executeErr).To(MatchError(shared.PluginNotFoundError{
					Name: "some-plugin",
				}))
			})
		})

		Context("when the actor returns a different error", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("some error")
				fakeSharedActor.UninstallPluginReturns(expectedErr)
			})

			It("returns the error", func() {
				Expect(executeErr).To(MatchError(expectedErr))
			})
		})
	})
})
