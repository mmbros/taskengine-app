package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/mmbros/flagx"
)

const usageApp = `Usage:
    %s <command> [options]
Available Commands:
    demo (d)     Execute a demo scenario, show progress
                 and output results in json format
    server (s)   Start an http server to show json files
                 containing the results of execution
    version (v)  Version information

Common options:
    -h, --help   Help informations
`

func initApp() *flagx.Command {

	app := &flagx.Command{
		ParseExec: parseExecApp,

		SubCmd: map[string]*flagx.Command{
			"version,v": {
				ParseExec: parseExecVersion,
			},
			"demo,d": {
				ParseExec: parseExecDemo,
			},
			"server,s": {
				ParseExec: parseExecServer,
			}},
	}

	return app
}

func parseExecApp(fullname string, arguments []string) error {
	fs := NewFlagSet(fullname, usageApp, fullname)
	err := fs.Parse(arguments)

	// handle help
	if err == nil || err == flag.ErrHelp {
		fs.Usage()
		return nil
	}

	return err
}

// Execute is the main function
func Execute(stderr io.Writer) int {
	app := initApp()
	if err := flagx.Run(app); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
