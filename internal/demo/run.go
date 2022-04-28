package demo

import (
	"encoding/json"
	"fmt"
	"io"
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

	etype := event.Type()

	jevent := struct {
		TaskID     string               `json:"task_id"`
		WorkerID   string               `json:"worker_id"`
		WorkerInst int                  `json:"worker_inst"`
		Status     taskengine.EventType `json:"status"`
		Label      string               `json:"label,omitempty"`
		TimeStart  time.Time            `json:"time_start"`
		TimeEnd    time.Time            `json:"time_end"`
		Err        string               `json:"err,omitempty"`
	}{
		WorkerID:   string(event.WorkerID),
		WorkerInst: event.WorkerInst,
		TaskID:     string(event.Task.TaskID()),
		TimeStart:  event.TimeStart,
		TimeEnd:    event.TimeEnd,
		Status:     etype,
		Label:      event.Result.String(),
	}

	if etype == taskengine.EventError {
		jevent.Err = event.Result.Error().Error()
	}

	//js, err := json.MarshalIndent(jevent, "", " ")
	js, err := json.Marshal(jevent)
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
		fmt.Fprint(wJson, "[\n")
	}

	var printComma bool
	for event := range eventc {

		// events = append(events, event)

		tid := string(event.Task.TaskID())

		progr.InitTrackerIfNew(tid)

		if wJson != nil && taskengine.IsResult(event) {
			if printComma {
				fmt.Fprint(wJson, ",\n")
			} else {
				printComma = true
			}

			fmt.Fprint(wJson, event2json(event))

		}

		if taskengine.IsFirstSuccessOrLastResult(event) {

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
		fmt.Fprint(wJson, "\n]")
	}

	progr.Render()
	time.Sleep(time.Millisecond * 200)

	progr.Stop()

	return &rs

}
