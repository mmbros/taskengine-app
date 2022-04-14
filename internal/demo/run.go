package demo

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/mmbros/taskengine"
)

func (scenario *Scenario) LoopWithProgress(eventc chan *taskengine.Event) {

	pw := NewProgress(scenario.Tasks)
	go pw.Render()
	trackers := map[string]*progress.Tracker{}

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
				tracker.Units = trackerUnits(result.currency)
				tracker.SetValue(trackerValue(result.price))
				tracker.UpdateMessage(tid + " " + text.FgHiBlack.Sprint(event.WorkerID))
				tracker.MarkAsDone()
			} else {
				tracker.UpdateMessage(tid + " " + text.FgRed.Sprint("error"))
				tracker.MarkAsErrored()
			}
		}

		// fmt.Println(event)
	}

	pw.Render()
	time.Sleep(time.Millisecond * 200)

	pw.Stop()

}
