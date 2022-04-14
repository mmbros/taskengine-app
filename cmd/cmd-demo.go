package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmbros/flagx"
	"github.com/mmbros/taskengine"
	"github.com/mmbros/taskengine-app/internal/demo"
)

const usageDemo = `Usage:
    %[1]s [options]
Options:
    -w, --workers     int      number of workers (default %[2]d)
    -i, --instances   int      instances of each worker (default %[3]d)
    -t, --tasks       int      number of tasks (default %[4]d)
        --progress    bool     show progress of execution (default %[5]v)
        --seed        int      xxxx (default %[6]d)
    -s, --spread      int      xxxx (default %[7]d)

Random Result options: 	
        --mean        int      mean value (default %[8]d)
        --stddev      int      standard deviation (default %[9]d)
        --errperc     int      perc of task error (0..100) (default %[10]d)
`

// stdDev:  100.0,
// mean:    500.0,
// errPerc: 50,

// Names of the command line arguments (flagx names)
const (
	namesWorkers   = "workers,w"
	namesInstances = "instances,i"
	namesTasks     = "tasks,t"
	namesProgress  = "progress"
	namesSeed      = "seed,d"
	namesSpread    = "spread"
	namesMean      = "mean"
	namesStdDev    = "stddev"
	namesErrPerc   = "errperc"
)

// Default args value
const (
	defaultWorkers   = 1
	defaultInstances = 1
	defaultTasks     = 10
	defaultProgress  = true
	defaultSeed      = 0
	defaultSpread    = 50
	defaultMean      = 500
	defaultStdDev    = 100
	defaultErrPerc   = 50
)

func parseExecDemo(fullname string, arguments []string) error {

	var (
		scenario     demo.Scenario
		showProgress bool
		seed         int
		mean         int
		stddev       int
	)

	fs := NewFlagSet(fullname, usageDemo,
		fullname, defaultWorkers, defaultInstances, defaultTasks, defaultProgress,
		defaultSeed, defaultSpread, defaultMean, defaultStdDev, defaultErrPerc)

	flagx.AliasedIntVar(fs, &scenario.Workers, namesWorkers, defaultWorkers, "")
	flagx.AliasedIntVar(fs, &scenario.Instances, namesInstances, defaultInstances, "")
	flagx.AliasedIntVar(fs, &scenario.Tasks, namesTasks, defaultTasks, "")
	flagx.AliasedBoolVar(fs, &showProgress, namesProgress, defaultProgress, "")
	flagx.AliasedIntVar(fs, &seed, namesSeed, defaultSeed, "")
	flagx.AliasedIntVar(fs, &scenario.Spread, namesSpread, defaultSpread, "")

	flagx.AliasedIntVar(fs, &mean, namesMean, defaultMean, "")
	flagx.AliasedIntVar(fs, &stddev, namesStdDev, defaultStdDev, "")
	flagx.AliasedIntVar(fs, &scenario.RandRes.ErrPerc, namesErrPerc, defaultErrPerc, "")

	// parse the arguments
	err := fs.Parse(arguments)

	// TODO flagx: create AliasedInt64Var, AliasedFloat64Var
	scenario.Seed = int64(seed)
	scenario.RandRes.Mean = float64(mean)
	scenario.RandRes.StdDev = float64(stddev)

	// handle help
	if err == flag.ErrHelp {
		fs.Usage()
		return nil
	}
	if err == nil {
		err = execDemo(os.Stdout, FirstToken(fullname, " "), &scenario)
	}
	return err
}

func execDemo(w io.Writer, appname string, scenario *demo.Scenario) error {

	var (
		err    error
		eventc chan *taskengine.Event
	)
	fmt.Fprintf(w, "%+v\n", scenario)

	err = scenario.RandomWorkersAndTasks()
	if err == nil {
		eventc, err = scenario.ExecuteEvents()
	}

	scenario.LoopWithProgress(eventc)

	if err != nil {
		return err
	}

	return err
}
