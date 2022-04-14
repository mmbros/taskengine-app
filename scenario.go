package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/mmbros/taskengine"
)

var demoError error = errors.New("taskengine demo error")

// spread: perc of how many workers executes each tasks:
//         100% - each task is executed by all worker
//           0% - no worker executes the tasks
type demoScenario struct {
	Seed      int64
	Workers   int
	Instances int
	Tasks     int
	Spread    int
	RandRes   randomResult

	ws  []*taskengine.Worker
	wts taskengine.WorkerTasks
}

// ========================================================

type randomResult struct {
	mean    float64
	stdDev  float64
	errPerc int
}

var rndPrice = randomResult{
	mean:   100,
	stdDev: 50,
}

func (rr *randomResult) float64() float64 {
	x := rand.NormFloat64()*rr.stdDev + rr.mean
	if x < 0 {
		x = 0
	}
	return x
}

func (rr *randomResult) int64() int64 {
	return int64(rr.float64())
}

func (rr *randomResult) success() bool {
	// errPerc = 0 .. 100
	// n       = 1 .. 100

	// if errPerc=  0 -> every n is greater than 0         -> always success
	// if errPerc=100 -> every n is less or equal than 100 -> always error
	n := rand.Intn(100) + 1
	return n > rr.errPerc
}

// ========================================================

type demoTask struct {
	taskid string
	rndres *randomResult
}

func (t *demoTask) TaskID() taskengine.TaskID { return taskengine.TaskID(t.taskid) }

// ========================================================

type demoResult struct {
	err      error
	price    float32
	currency string
}

func (res *demoResult) Error() error { return res.err }

func (res *demoResult) TrackerUnits() progress.Units {
	u := UnitsCurrency
	var currency string
	switch res.currency {
	case "USD":
		currency = "$"
	case "GBP":
		currency = "£"
	case "EUR":
		currency = "€"
	default:
		currency = res.currency
	}
	u.Notation = currency + " "

	return u
}

func (res *demoResult) TrackerValue() int64 {
	return int64(res.price * 100)
}

// ========================================================

func demoWorkFn(ctx context.Context, worker *taskengine.Worker, workerInst int, task taskengine.Task) taskengine.Result {
	stask := task.(*demoTask)

	msec := stask.rndres.int64()

	res := &demoResult{}

	select {
	case <-ctx.Done():
		res.err = ctx.Err()
	case <-time.After(time.Duration(msec) * time.Millisecond):
		if !stask.rndres.success() {
			res.err = demoError
		} else {
			res.price = float32(rndPrice.float64())
			if msec%2 == 0 {
				res.currency = "EUR"
			} else if msec%3 == 0 {
				res.currency = "USD"
			} else {
				res.currency = "GBP"
			}
		}
	}

	return taskengine.Result(res)
}

func getDigits(tot int) int {
	var digits int
	if tot >= 1000 {
		digits = 4
	} else if tot >= 100 {
		digits = 3
	} else if tot >= 10 {
		digits = 2
	} else {
		digits = 1
	}
	return digits
}
func getID(prefix string, tot, n int) string {
	digits := getDigits(tot)
	return fmt.Sprintf("%[3]s%0[2]*[1]d", n, digits, prefix)
}

func (scnr *demoScenario) getWorkerID(n int) string {
	return getID("w", scnr.Workers, n)
}

func (scnr *demoScenario) getTaskID(n int) string {
	return getID("t", scnr.Tasks, n)
}

func (scnr *demoScenario) RandomWorkersAndTasks() error {

	rand.Seed(scnr.Seed)
	ws := []*taskengine.Worker{}
	wts := taskengine.WorkerTasks{}
	if scnr.Spread < 0 || scnr.Spread > 100 {
		return errors.New("Spread must be in 0..100")
	}

	for wj := 1; wj <= scnr.Workers; wj++ {
		wid := taskengine.WorkerID(scnr.getWorkerID(wj))
		w := &taskengine.Worker{
			WorkerID:  wid,
			Instances: scnr.Instances,
			Work:      demoWorkFn,
		}
		ws = append(ws, w)

		ts := taskengine.Tasks{}
		for tj := 1; tj <= scnr.Tasks; tj++ {

			n := rand.Intn(101)
			if n <= scnr.Spread {
				// assign the task to the worker
				tid := scnr.getTaskID(tj)
				task := &demoTask{
					taskid: tid,
					rndres: &scnr.RandRes,
				}
				ts = append(ts, task)
			}
		}
		wts[wid] = ts
	}

	//  add tasks not used
	t2w := GetTaskWorkers(wts)
	for tj := 1; tj <= scnr.Tasks; tj++ {
		tid := scnr.getTaskID(tj)
		if _, ok := t2w[taskengine.TaskID(tid)]; !ok {
			// the task has no workers:
			// assign the task to a worker
			wj := rand.Intn(scnr.Workers) + 1
			wid := taskengine.WorkerID(scnr.getWorkerID(wj))
			task := &demoTask{
				taskid: tid,
				rndres: &scnr.RandRes,
			}
			ts, ok := wts[wid]
			if !ok {
				ts = taskengine.Tasks{task}
			} else {
				ts = append(ts, task)
			}
			wts[wid] = ts
		}
	}

	scnr.ws = ws
	scnr.wts = wts

	return nil
}

func (scnr *demoScenario) ExecuteEvents() (chan *taskengine.Event, error) {

	if scnr.ws == nil {
		return nil, errors.New("must run RandomWorkersAndTasks before")
	}

	ctx := context.Background()
	eng, err := taskengine.NewEngine(scnr.ws, scnr.wts)
	if err != nil {
		return nil, err
	}
	return eng.ExecuteEvents(ctx)
}
