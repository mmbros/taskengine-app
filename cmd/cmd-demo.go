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

// Spread: perc of how many workers executes each tasks
// 100% - each task is executed by all worker
//   0% - no worker executes the tasks

const usageDemo = `Usage:
    %[1]s [options]

Performs a demo scenario and show results. 

Options:
    -w, --workers     int      number of workers (default %[2]d)
    -i, --instances   int      instances of each worker (default %[3]d)
    -t, --tasks       int      number of tasks (default %[4]d)
        --progress    bool     show progress of execution (default %[5]v)
        --seed        int      random seed generator (default %[6]d)
    -s, --spread      int      perc of how many workers executes each tasks (default %[7]d)
                                 100 each task is executed by all worker
                                   0 no worker executes the tasks
    -o, --output      path     pathname of the output file (default stdout)
    -f, --force       bool     overwrite already existing output file

Random Result options: 	
        --mean        int      mean value (default %[8]d)
        --stddev      int      standard deviation (default %[9]d)
    -e, --errperc     int      perc of task error (0..100) (default %[10]d)

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
	namesErrPerc   = "errperc,e"
	namesOutput    = "output,o"
	namesForce     = "force,f"
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
		scenario demo.Scenario
		seed     int
		mean     int
		stddev   int

		path         string
		showProgress bool
		force        bool
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
	flagx.AliasedStringVar(fs, &path, namesOutput, "", "")
	flagx.AliasedBoolVar(fs, &force, namesForce, false, "")

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

	wOutput := os.Stdout
	if path != "" {
		// overwrite existing file oly if --force is specified
		flag := os.O_CREATE | os.O_WRONLY
		if !force {
			flag |= os.O_EXCL
		}

		// create the file
		wOutput, err = os.OpenFile(path, flag, 0666)
		if err == nil {
			defer wOutput.Close()
		}
	}

	if err == nil {
		err = execDemo(os.Stderr, wOutput, &scenario, showProgress)
	}
	return err
}

func execDemo(wInfo, wOut io.Writer, scenario *demo.Scenario, showProgress bool) error {

	var (
		err    error
		eventc chan *taskengine.Event
	)
	if wInfo != nil {
		fmt.Fprintf(wInfo, "%+v\n", scenario)
	}

	err = scenario.RandomWorkersAndTasks()

	if err == nil {
		eventc, err = scenario.ExecuteEvents()
	}
	if err != nil {
		return err
	}

	var wInfo2 io.Writer
	if wInfo != nil && showProgress {
		wInfo2 = wInfo
	}
	stats := scenario.Run(eventc, wInfo2, wOut)

	if wInfo != nil {
		fmt.Fprintf(wInfo, "\n%d task completed (%d success, %d error) in %v\n",
			stats.TaskCompleted(),
			stats.TaskSuccess,
			stats.TaskError,
			stats.Elapsed())
	}

	return err
}
