package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmbros/flagx"
)

const usageVersion = `Usage:
    %s [-b | --build-options]

Prints version informations.

Options:
    -b, --build-options  bool   print verion with build options
`

// Names of the command line arguments (flagx names)
const (
	namesBuildOptions = "build-options,b"
)

// Default args value
const (
	defaultBuildOptions = false
)

// set at compile time with
//   -ldflags="-X 'github.com/mmbros/quotes/cmd.AppVersion=x.y.z' -X 'github.com/mmbros/quotes/cmd.GitCommit=...'"
var (
	AppVersion        string // git tag ...
	TaskEngineVersion string // git tag ...
	GitCommit         string // git rev-parse --short HEAD
	GoVersion         string // go version
	BuildTime         string // when the executable was built
	OsArch            string // uname -s -m
)

func parseExecVersion(fullname string, arguments []string) error {
	var buildOptions bool

	fs := NewFlagSet(fullname, usageVersion, fullname)

	flagx.AliasedBoolVar(fs, &buildOptions, namesBuildOptions, defaultBuildOptions, "")

	// parse the arguments
	err := fs.Parse(arguments)

	// handle help
	if err == flag.ErrHelp {
		fs.Usage()
		return nil
	}
	if err == nil {
		execVersion(os.Stdout, FirstToken(fullname, " "), buildOptions)
	}
	return err
}

func execVersion(w io.Writer, appname string, buildOptions bool) {
	fmt.Fprintf(w, "%s version %s\ntaskengine package version %s\n", appname, AppVersion, TaskEngineVersion)
	if buildOptions {
		fmt.Fprintf(w, `%s
build date: %s
git commit: %s
os/arch: %s
`, GoVersion, BuildTime, GitCommit, OsArch)
	}
}
