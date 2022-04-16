package cmd

import (
	"flag"

	"github.com/mmbros/flagx"
	"github.com/mmbros/taskengine-app/internal/server"
)

const usageView = `Usage:
    %[1]s
Start an http server to read the json files and show a graph of workers executions.

Options:
    -f, --folder        folder containing the json files (default %[2]q)
    -a, --address       server address and port (default %[3]q)
`

// Names of the command line arguments (flagx names)
const (
	namesFolder  = "folder,f"
	namesAddress = "address,a"
)

// Default args value
const (
	defaultFolder  = "."
	defaultAddress = ":6789"
)

func parseExecServer(fullname string, arguments []string) error {
	var folder string
	var address string

	fs := NewFlagSet(fullname, usageView,
		fullname, defaultFolder, defaultAddress)

	flagx.AliasedStringVar(fs, &folder, namesFolder, defaultFolder, "")
	flagx.AliasedStringVar(fs, &address, namesAddress, defaultAddress, "")

	// parse the arguments
	err := fs.Parse(arguments)

	// handle help
	if err == flag.ErrHelp {
		fs.Usage()
		return nil
	}
	return execServer(address, folder)
}

func execServer(serverAddressPort, jsonDataFolder string) error {
	return server.Run(serverAddressPort, jsonDataFolder)
}
