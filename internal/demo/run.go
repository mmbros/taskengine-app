package demo

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mmbros/taskengine"
	"github.com/mmbros/taskengine-app/internal/progress"
)

type RunStats struct {
	TaskSuccess int
	TaskError   int

	TimeStart time.Time
	TimeEnd   time.Time
}

func (stats *RunStats) Elapsed() time.Duration {
	return stats.TimeEnd.Sub(stats.TimeStart)
}
func (stats *RunStats) TaskCompleted() int {
	return stats.TaskSuccess + stats.TaskError
}

// ========================================================

func event2json(event *taskengine.Event) string {

	jevent := struct {
		WorkerID   string    `json:"worker_id"`
		WorkerInst int       `json:"worker_inst"`
		TaskID     string    `json:"task_id"`
		TimeStart  time.Time `json:"time_start"`
		TimeEnd    time.Time `json:"time_end"`
		Status     string    `json:"status"`
	}{
		WorkerID:   string(event.WorkerID),
		WorkerInst: event.WorkerInst,
		TaskID:     string(event.Task.TaskID()),
		TimeStart:  event.TimeStart,
		TimeEnd:    event.TimeEnd,
		Status:     strings.ToLower(event.Type().String()),
	}

	js, err := json.MarshalIndent(jevent, "", " ")
	// js, err := json.Marshal(jevent)
	if err != nil {
		panic(err)
	}
	return string(js)
}

func (scenario *Scenario) Run(eventc chan *taskengine.Event, wProgress, wJson io.Writer) *RunStats {

	var progr *progress.Progress
	if wProgress != nil {
		progr = progress.New(wProgress, scenario.Tasks)
	}

	go progr.Render()

	rs := RunStats{}
	rs.TimeStart = time.Now()

	if wJson != nil {
		fmt.Fprint(wJson, "[")
	}

	var printComma bool
	for event := range eventc {

		// events = append(events, event)

		tid := string(event.Task.TaskID())

		progr.InitTrackerIfNew(tid)

		if wJson != nil && event.IsResult() {
			if printComma {
				fmt.Fprint(wJson, ",")
			} else {
				printComma = true
			}

			fmt.Fprint(wJson, event2json(event))

		}

		if event.IsFirstSuccessOrLastResult() {

			if event.Type() == taskengine.EventSuccess {
				result := event.Result.(*demoResult)
				progr.SetSuccess(tid, string(event.WorkerID), result.price, result.currency)
				rs.TaskSuccess++
			} else {
				progr.SetError(tid)
				rs.TaskError++
			}
		}

		// fmt.Println(event)
	}

	rs.TimeEnd = time.Now()

	if wJson != nil {
		fmt.Fprint(wJson, "]")
	}

	progr.Render()
	time.Sleep(time.Millisecond * 200)

	progr.Stop()

	return &rs

}
