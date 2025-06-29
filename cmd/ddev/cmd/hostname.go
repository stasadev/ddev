package cmd

import (
	"github.com/ddev/ddev/pkg/ddevapp"
	"github.com/ddev/ddev/pkg/exec"
	"github.com/ddev/ddev/pkg/hostname"
	"github.com/ddev/ddev/pkg/util"
	"github.com/spf13/cobra"
	"os"
)

var removeHostnameFlag bool
var removeInactiveFlag bool
var checkHostnameFlag bool

// HostNameCmd represents the hostname command
var HostNameCmd = &cobra.Command{
	Use:   "hostname [hostname] [ip]",
	Short: "Manage your hostfile entries.",
	Example: `
ddev hostname junk.example.com 127.0.0.1
ddev hostname -r junk.example.com 127.0.0.1
ddev hostname --check junk.example.com 127.0.0.1
ddev hostname --remove-inactive
`,
	Long: `Manage your hostfile entries. Managing host names has security and usability
implications and requires elevated privileges. You may be asked for a password
to allow DDEV to modify your hosts file. If you are connected to the internet and using the domain ddev.site this is generally not necessary, because the hosts file never gets manipulated.`,
	Run: func(_ *cobra.Command, args []string) {
		ddevHostnameBinary := hostname.GetDdevHostnameBinary()
		// If requested, remove all inactive host names and exit
		if removeInactiveFlag {
			if len(args) > 0 {
				util.Failed("Invalid arguments supplied. 'ddev hostname --remove-all' accepts no arguments.")
			}

			util.Warning("Attempting to remove inactive custom hostnames for projects which are registered but not running")
			removeInactiveHostnames()
			return
		}

		// If operating on one host name, two arguments are required
		if len(args) != 2 {
			util.Failed("Invalid arguments supplied. Please use 'ddev hostname [hostname] [ip]'")
		}

		name, dockerIP := args[0], args[1]
		var err error

		hostnameInHostsFile, err := hostname.IsHostnameInHostsFile(name)
		if err != nil {
			util.Warning("Could not check existence of %s in hosts file: %v", name, err)
		}
		// If requested, remove the provided host name and exit
		if removeHostnameFlag {
			if !hostnameInHostsFile {
				return
			}
			out, err := exec.RunHostCommand(ddevHostnameBinary, "--remove", name, dockerIP)
			if err != nil {
				util.Warning("Failed to remove hosts entry %s:%s: %v (output='%v')", name, dockerIP, err, out)
			}
			return
		}
		if checkHostnameFlag {
			if hostnameInHostsFile {
				return
			}

			os.Exit(1)
		}
		// By default, add a host name
		if !hostnameInHostsFile {
			out, err := exec.RunHostCommand(ddevHostnameBinary, name, dockerIP)
			if err != nil {
				util.Warning("Failed to add hosts entry %s:%s: %v (output='%v')", name, dockerIP, err, out)
			}
		}
	},
}

// removeInactiveHostnames will remove all host names except those current in use by active projects.
func removeInactiveHostnames() {
	apps, err := ddevapp.GetInactiveProjects()
	if err != nil {
		util.Warning("Unable to run GetInactiveProjects: %v", err)
		return
	}
	for _, app := range apps {
		err := app.RemoveHostsEntriesIfNeeded()
		if err != nil {
			util.Warning("Unable to remove hosts entries for project '%s': %v", app.Name, err)
		}
	}
}

func init() {
	HostNameCmd.Flags().BoolVarP(&removeHostnameFlag, "remove", "r", false, "Remove the provided host name - ip correlation")
	HostNameCmd.Flags().BoolVarP(&checkHostnameFlag, "check", "c", false, "Check to see if provided hostname is already in hosts file")
	HostNameCmd.Flags().BoolVarP(&removeInactiveFlag, "remove-inactive", "R", false, "Remove host names of inactive projects")
	HostNameCmd.Flags().BoolVar(&removeInactiveFlag, "fire-bazooka", false, "Alias of --remove-inactive")
	_ = HostNameCmd.Flags().MarkHidden("fire-bazooka")

	RootCmd.AddCommand(HostNameCmd)
}
