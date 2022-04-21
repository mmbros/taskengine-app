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

Performs a demo scenario and save results in json format. 

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

Examples:
    # print results to file with progress
    %[1]s --output out.json

    # print results to stdout without progress
    %[1]s --progress=0
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

type params struct {
	force    bool
	output   string
	progress bool
	scenario demo.Scenario
}

func (p *params) String() string {
	// return fmt.Sprintf("%#v", p)
	return fmt.Sprintf("force=%v output=%q progress=%v seed=%d workers=%d instances=%d tasks=%d spread=%v mean=%v stddev=%v errperc=%d",
		p.force, p.output, p.progress, p.scenario.Seed,
		p.scenario.Workers, p.scenario.Instances, p.scenario.Tasks,
		p.scenario.Spread,
		p.scenario.RandRes.Mean, p.scenario.RandRes.StdDev, p.scenario.RandRes.ErrPerc,
	)
}

func parseExecDemo(fullname string, arguments []string) error {

	var p params

	fs := NewFlagSet(fullname, usageDemo,
		fullname, defaultWorkers, defaultInstances, defaultTasks, defaultProgress,
		defaultSeed, defaultSpread, defaultMean, defaultStdDev, defaultErrPerc)

	flagx.AliasedBoolVar(fs, &p.force, namesForce, false, "")
	flagx.AliasedStringVar(fs, &p.output, namesOutput, "", "")
	flagx.AliasedBoolVar(fs, &p.progress, namesProgress, defaultProgress, "")

	flagx.AliasedIntVar(fs, &p.scenario.Workers, namesWorkers, defaultWorkers, "")
	flagx.AliasedIntVar(fs, &p.scenario.Instances, namesInstances, defaultInstances, "")
	flagx.AliasedIntVar(fs, &p.scenario.Tasks, namesTasks, defaultTasks, "")
	flagx.AliasedInt64Var(fs, &p.scenario.Seed, namesSeed, defaultSeed, "")
	flagx.AliasedIntVar(fs, &p.scenario.Spread, namesSpread, defaultSpread, "")

	flagx.AliasedFloat64Var(fs, &p.scenario.RandRes.Mean, namesMean, defaultMean, "")
	flagx.AliasedFloat64Var(fs, &p.scenario.RandRes.StdDev, namesStdDev, defaultStdDev, "")
	flagx.AliasedIntVar(fs, &p.scenario.RandRes.ErrPerc, namesErrPerc, defaultErrPerc, "")

	// parse the arguments
	err := fs.Parse(arguments)

	// handle help
	if err == flag.ErrHelp {
		fs.Usage()
		return nil
	}

	if err != nil {
		return err
	}

	wOutput := os.Stdout
	if p.output != "" {
		// overwrite existing file only if --force is specified
		flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		if !p.force {
			flag |= os.O_EXCL // file must not exists
		}

		// create the file
		wOutput, err = os.OpenFile(p.output, flag, 0666)
		if err != nil {
			return err
		}
		defer wOutput.Close()
	}

	err = execDemo(os.Stderr, wOutput, &p)

	return err
}

func execDemo(wInfo, wOut io.Writer, p *params) error {

	var (
		err    error
		eventc chan *taskengine.Event
	)
	if wInfo != nil {
		fmt.Fprintln(wInfo, p.String())
	}

	err = p.scenario.RandomWorkersAndTasks()

	if err == nil {
		eventc, err = p.scenario.ExecuteEvents()
	}
	if err != nil {
		return err
	}

	var wProgress io.Writer
	if wInfo != nil && p.progress {
		wProgress = wInfo
	}
	stats := p.scenario.Run(eventc, wProgress, wOut)

	if wInfo != nil {
		fmt.Fprintf(wInfo, "\n%d task completed (%d success, %d error) in %v\n",
			stats.TaskCompleted(),
			stats.TaskSuccess,
			stats.TaskError,
			stats.Elapsed())
	}

	return err
}
