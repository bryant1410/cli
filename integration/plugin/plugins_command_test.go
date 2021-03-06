package plugin

import (
	"os"
	"path/filepath"

	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("plugins command", func() {
	BeforeEach(func() {
		helpers.RunIfExperimental("Running experimental plugins command tests")
		// This removes plugin artefacts from other plugin tests
		uninstallTestPlugin()
	})

	Describe("help", func() {
		Context("when --help flag is provided", func() {
			It("displays command usage to output", func() {
				session := helpers.CF("plugins", "--help")
				Eventually(session).Should(Say("NAME:"))
				Eventually(session).Should(Say("plugins - List all available plugin commands"))
				Eventually(session).Should(Say("USAGE:"))
				Eventually(session).Should(Say("cf plugins [--checksum]"))
				Eventually(session).Should(Say("OPTIONS:"))
				Eventually(session).Should(Say("--checksum\\s+Compute and show the sha1 value of the plugin binary file"))
				Eventually(session).Should(Say("SEE ALSO:"))
				Eventually(session).Should(Say("install-plugin, repo-plugins, uninstall-plugin"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when no plugins are installed", func() {
		It("displays an empty table", func() {
			session := helpers.CF("plugins")
			Eventually(session).Should(Say("Listing installed plugins..."))
			Eventually(session).Should(Say(""))
			Eventually(session).Should(Say("plugin name\\s+version\\s+command name\\s+command help"))
			Consistently(session).ShouldNot(Say("[a-za-z0-9]+"))
			Eventually(session).Should(Exit(0))
		})

		Context("when the --checksum flag is provided", func() {
			It("displays an empty checksum table", func() {
				session := helpers.CF("plugins", "--checksum")
				Eventually(session).Should(Say("Computing sha1 for installed plugins, this may take a while..."))
				Eventually(session).Should(Say(""))
				Eventually(session).Should(Say("plugin name\\s+version\\s+sha1"))
				Consistently(session).ShouldNot(Say("[a-za-z0-9]+"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when plugins are installed", func() {
		Context("when there are multiple plugins", func() {
			BeforeEach(func() {
				helpers.CreateBasicPlugin("I-should-be-sorted-first", "1.2.0", []helpers.PluginCommand{
					{Name: "command-1", Help: "some-command-1"},
					{Name: "Better-command", Help: "some-better-command"},
					{Name: "command-2", Help: "some-command-2"},
				})
				helpers.CreateBasicPlugin("sorted-third", "2.0.1", []helpers.PluginCommand{
					{Name: "banana-command", Help: "banana-command"},
				})
				helpers.CreateBasicPlugin("i-should-be-sorted-second", "1.0.0", []helpers.PluginCommand{
					{Name: "some-command", Help: "some-command"},
					{Name: "Some-other-command", Help: "some-other-command"},
				})
			})

			It("displays the installed plugins in alphabetical order", func() {
				session := helpers.CF("plugins")
				Eventually(session).Should(Say("Listing installed plugins..."))
				Eventually(session).Should(Say(""))
				Eventually(session).Should(Say("plugin name\\s+version\\s+command name\\s+command help"))
				Eventually(session).Should(Say("I-should-be-sorted-first\\s+1\\.2\\.0\\s+command-1\\s+some-command-1"))
				Eventually(session).Should(Say("I-should-be-sorted-first\\s+1\\.2\\.0\\s+Better-command\\s+some-better-command"))
				Eventually(session).Should(Say("I-should-be-sorted-first\\s+1\\.2\\.0\\s+command-2\\s+some-command-2"))
				Eventually(session).Should(Say("i-should-be-sorted-second\\s+1\\.0\\.0\\s+some-command\\s+some-command"))
				Eventually(session).Should(Say("i-should-be-sorted-second\\s+1\\.0\\.0\\s+Some-other-command\\s+some-other-command"))
				Eventually(session).Should(Say("sorted-third\\s+2\\.0\\.1\\s+banana-command\\s+banana-command"))
				Eventually(session).Should(Exit(0))
			})
		})

		Context("when plugin version information is 0.0.0", func() {
			BeforeEach(func() {
				helpers.CreateBasicPlugin("some-plugin", "0.0.0", []helpers.PluginCommand{
					{Name: "banana-command", Help: "banana-command"},
				})
			})

			It("displays N/A for the plugin's version", func() {
				session := helpers.CF("plugins")
				Eventually(session).Should(Say("some-plugin\\s+N/A"))
				Eventually(session).Should(Exit(0))
			})
		})

		Context("when a plugin command has an alias", func() {
			BeforeEach(func() {
				helpers.CreateBasicPlugin("some-plugin", "1.0.0", []helpers.PluginCommand{
					{Name: "banana-command", Alias: "bc", Help: "banana-command"},
				})
			})

			It("displays the command name and it's alias", func() {
				session := helpers.CF("plugins")
				Eventually(session).Should(Say("some-plugin\\s+1\\.0\\.0\\s+banana-command, bc"))
				Eventually(session).Should(Exit(0))
			})
		})

		Context("when the --checksum flag is provided", func() {
			BeforeEach(func() {
				helpers.CreateBasicPlugin("some-plugin", "1.0.0", []helpers.PluginCommand{
					{Name: "banana-command", Help: "banana-command"},
				})
			})

			It("displays the sha1 value for each installed plugin", func() {
				session := helpers.CF("plugins", "--checksum")
				Eventually(session).Should(Say("Computing sha1 for installed plugins, this may take a while..."))
				Eventually(session).Should(Say(""))
				Eventually(session).Should(Say("plugin name\\s+version\\s+sha1"))
				Eventually(session).Should(Say("some-plugin\\s+1\\.0\\.0\\s+ee3f5d2802f8b332fae9589fc367537ff75a41b0"))
				Eventually(session).Should(Exit(0))
			})

			Context("when an error is encountered calculating the sha1 value", func() {
				It("displays N/A for the plugin's sha1", func() {
					err := os.Remove(filepath.Join(homeDir, ".cf/plugins/configurable_plugin.some-plugin"))
					Expect(err).NotTo(HaveOccurred())

					session := helpers.CF("plugins", "--checksum")
					Eventually(session).Should(Say("some-plugin\\s+1\\.0\\.0\\s+N/A"))
					Eventually(session).Should(Exit(0))
				})
			})
		})

	})
})
