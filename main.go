package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/text"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/mmbros/taskengine"
)

var (
	flagNoProgress = flag.Bool("no-progress", false, "Do not show progress")
	flagSeed       = flag.Int64("seed", 0, "Seed of the rand function")
	flagWorkers    = flag.Int("workers", 3, "Number of workers")
	flagInstances  = flag.Int("instances", 2, "Number of workers instances")
	flagTasks      = flag.Int("tasks", 20, "Number of tasks")
	flagSpread     = flag.Int("spread", 100, "Perc. of how many workers executes each tasks")
)

type runResult struct {
	numTaskSuccess int
	numTaskError   int
	timeStart      time.Time
	timeEnd        time.Time
}

func runScenarioWithProgress(scenario *demoScenario, eventc chan *taskengine.Event) *runResult {
	rr := &runResult{}

	pw := initProgress(scenario.Tasks)
	go pw.Render()
	trackers := map[string]*progress.Tracker{}

	rr.timeStart = time.Now()

	for event := range eventc {

		tid := string(event.Task.TaskID())

		tracker, ok := trackers[tid]
		if !ok {
			tracker = &progress.Tracker{
				Message: tid,
				Total:   0,
				Units:   UnitsNA,
			}
			pw.AppendTracker(tracker)
			trackers[tid] = tracker
		}

		if event.IsFirstSuccessOrLastResult() {

			if event.Type() == taskengine.EventSuccess {
				result := event.Result.(*demoResult)
				tracker.Units = result.TrackerUnits()
				tracker.SetValue(result.TrackerValue())
				tracker.UpdateMessage(tid + " " + text.FgHiBlack.Sprint(event.WorkerID))
				tracker.MarkAsDone()
				rr.numTaskSuccess++
			} else {
				tracker.UpdateMessage(tid + " " + text.FgRed.Sprint("error"))
				tracker.MarkAsErrored()
				rr.numTaskError++
			}
		}

		// fmt.Println(event)
	}
	rr.timeEnd = time.Now()

	pw.Render()
	time.Sleep(time.Millisecond * 200)

	pw.Stop()

	return rr
}

func runScenario(scenario *demoScenario, eventc chan *taskengine.Event) *runResult {
	rr := &runResult{}

	rr.timeStart = time.Now()

	for event := range eventc {
		if event.IsFirstSuccessOrLastResult() {
			if event.Type() == taskengine.EventSuccess {
				rr.numTaskSuccess++
			} else {
				rr.numTaskError++
			}
		}
		// fmt.Println(event)
	}
	rr.timeEnd = time.Now()
	return rr
}

func main() {

	var (
		eventc chan *taskengine.Event
		err    error
		rr     *runResult
	)

	flag.Parse()

	scenario := demoScenario{
		Seed:      *flagSeed,
		Workers:   *flagWorkers,
		Instances: *flagInstances,
		Tasks:     *flagTasks,
		Spread:    *flagSpread,
		RandRes: randomResult{
			stdDev:  100.0,
			mean:    500.0,
			errPerc: 50,
		},
	}

	err = scenario.RandomWorkersAndTasks()
	if err == nil {
		eventc, err = scenario.ExecuteEvents()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("\nStart %d tasks with %d workers x %d instances\n",
		scenario.Tasks,
		scenario.Workers,
		scenario.Instances,
	)

	if *flagNoProgress {
		rr = runScenario(&scenario, eventc)
	} else {
		rr = runScenarioWithProgress(&scenario, eventc)
	}

	elapsed := rr.timeEnd.Sub(rr.timeStart).Seconds()

	fmt.Printf("\nFinished %d tasks (%d success, %d error) in %.3fs with %d workers x %d instances\n",
		rr.numTaskSuccess+rr.numTaskError,
		rr.numTaskSuccess,
		rr.numTaskError,
		elapsed,
		len(scenario.ws),
		scenario.Instances,
	)

}
