package cmd

import (
	"flag"
	"fmt"
	"strings"
)

func FirstToken(s string, sep string) string {
	astr := strings.Split(s, sep)
	if len(astr) == 0 {
		return ""
	}
	return astr[0]
}

func NewFlagSet(fullname string, usage string, a ...interface{}) *flag.FlagSet {

	fs := flag.NewFlagSet(fullname, flag.ContinueOnError)

	// use the same output as flag.CommandLine
	fs.SetOutput(flag.CommandLine.Output())

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), usage, a...)
	}

	return fs
}

// import (
// 	"flag"
// 	"fmt"
// 	"io"
// 	"strings"

// 	"github.com/mmbros/flagx"
// )

// type flagsGroup int

// // Type of flags group

// const (
// 	fgApp flagsGroup = iota
// 	fgAppVersion
// )

// // Names of the command line arguments (flagx names)
// const (
// 	namesWorkers    = "workers,w"
// 	namesInstanced  = "instances,i"
// 	namesTasks      = "tasks,t"
// 	namesNoProgress = "no-progress"
// 	namesSeed       = "seed,d"
// 	namesSpread     = "spread,p"
// )

// // Default args value
// const (
// 	defaultWorkers   = 1
// 	defaultInstances = 1
// 	defaultTasks     = 10
// 	defaultSeed      = 0
// 	defaultSpread    = 50
// )

// type Flags struct {
// 	workers    int
// 	instances  int
// 	tasks      int
// 	noProgress bool
// 	seed       int
// 	spread     int

// 	flagSet  *flag.FlagSet
// 	fullname string
// }

// func NewFlags(fullname string, flagsgroup flagsGroup) *Flags {
// 	fs := flag.NewFlagSet(fullname, flag.ContinueOnError)

// 	flags := &Flags{}
// 	flags.flagSet = fs
// 	flags.fullname = fullname

// 	// use the same output as flag.CommandLine
// 	fs.SetOutput(flag.CommandLine.Output())

// 	// flagx.AliasedBoolVar(fs, &flags.dryrun, namesDryrun, false, "")
// 	// flagx.AliasedIntVar(fs, &flags.workers, namesWorkers, defaultWorkers, "")
// 	// flagx.AliasedStringVar(fs, &flags.database, namesDatabase, "", "")
// 	// flagx.AliasedStringVar(fs, &flags.mode, namesMode, defaultMode, "")
// 	// flagx.AliasedStringsVar(fs, &flags.isins, namesIsins, "")
// 	// flagx.AliasedStringsVar(fs, &flags.sources, namesSources, "")

// 	return flags
// }

// // SetUsage set the usage function of the inner FlagSet
// func (flags *Flags) SetUsage(format string, a ...interface{}) {
// 	fs := flags.flagSet
// 	fs.Usage = func() {
// 		fmt.Fprintf(fs.Output(), format, a...)
// 	}
// }

// // IsPassed checks if the flag was passed in the command-line arguments.
// // names is a string that contains the comma separated aliases of the flag.
// func (flags *Flags) IsPassed(names string) bool {
// 	return flagx.IsPassed(flags.flagSet, names)
// }

// // Appname returns the app name from the fullname of the command
// //
// // Example:
// //   fullname = "app cmd sub-cmd sub-sub-cmd"
// //   output   = "app"
// func (flags *Flags) Appname() string {
// 	astr := strings.Split(flags.fullname, " ")
// 	if len(astr) == 0 {
// 		return ""
// 	}
// 	return astr[0]
// }

// func (flags *Flags) Parse(arguments []string) error { return flags.flagSet.Parse(arguments) }

// func (flags *Flags) Usage() { flags.flagSet.Usage() }

// func (flags *Flags) Output() io.Writer { return flags.flagSet.Output() }
