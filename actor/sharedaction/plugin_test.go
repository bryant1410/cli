package sharedaction_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/sharedaction/sharedactionfakes"
	"code.cloudfoundry.org/cli/util/configv3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UninstallPlugin", func() {
	var (
		actor                 Actor
		binaryName            string
		fakeConfig            *sharedactionfakes.FakeConfig
		fakePluginUninstaller *sharedactionfakes.FakePluginUninstaller
		pluginHome            string
		err                   error
	)

	BeforeEach(func() {
		pluginHome, err = ioutil.TempDir("", "")
		Expect(err).ToNot(HaveOccurred())

		binaryName = filepath.Join(pluginHome, "banana-faceman")
		err = ioutil.WriteFile(binaryName, nil, 0600)
		Expect(err).ToNot(HaveOccurred())

		fakeConfig = new(sharedactionfakes.FakeConfig)
		plugins := map[string]configv3.Plugin{
			"some-plugin": configv3.Plugin{
				Location: binaryName,
			},
		}
		fakeConfig.PluginsReturns(plugins)
		fakePluginUninstaller = new(sharedactionfakes.FakePluginUninstaller)

		actor = NewActor()
	})

	AfterEach(func() {
		os.RemoveAll(pluginHome)
	})

	Context("when the plugin does not exist", func() {
		It("returns a PluginNotFoundError", func() {
			err := actor.UninstallPlugin(fakeConfig, fakePluginUninstaller, "some-non-existant-plugin")
			Expect(err).To(MatchError(PluginNotFoundError{Name: "some-non-existant-plugin"}))
		})
	})

	Context("when the plugin exists", func() {
		It("runs the plugin cleanup, deletes the binary and removes the plugin config", func() {
			err := actor.UninstallPlugin(fakeConfig, fakePluginUninstaller, "some-plugin")
			Expect(err).ToNot(HaveOccurred())

			Expect(fakePluginUninstaller.UninstallCallCount()).To(Equal(1))
			Expect(fakePluginUninstaller.UninstallArgsForCall(0)).To(Equal(binaryName))

			_, err = os.Stat(binaryName)
			Expect(os.IsNotExist(err)).To(BeTrue())

			Expect(fakeConfig.PluginsCallCount()).To(Equal(1))

			Expect(fakeConfig.RemovePluginCallCount()).To(Equal(1))
			Expect(fakeConfig.RemovePluginArgsForCall(0)).To(Equal("some-plugin"))
		})

		Context("when the plugin uninstaller returns an error", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("some error")
				fakePluginUninstaller.UninstallReturns(expectedErr)
			})

			It("returns the error and does not delete the binary or remove the plugin config", func() {
				err := actor.UninstallPlugin(fakeConfig, fakePluginUninstaller, "some-plugin")
				Expect(err).To(MatchError(expectedErr))

				_, err = os.Stat(binaryName)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeConfig.RemovePluginCallCount()).To(Equal(0))
			})
		})

		Context("when deleting the plugin binary returns an error", func() {
			BeforeEach(func() {
				Expect(os.Remove(binaryName)).ToNot(HaveOccurred())
			})

			It("returns an error and does not remove the plugin config", func() {
				err := actor.UninstallPlugin(fakeConfig, fakePluginUninstaller, "some-plugin")
				Expect(err).To(MatchError(MatchRegexp("no such file or directory")))

				Expect(fakeConfig.RemovePluginCallCount()).To(Equal(0))
			})
		})
	})
})
