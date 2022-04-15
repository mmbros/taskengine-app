package demo

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/mmbros/taskengine"
)

var ErrorResult error = errors.New("demo error result")

type Scenario struct {
	Seed      int64
	Workers   int
	Instances int
	Tasks     int

	// Spread: perc of how many workers executes each tasks
	// 100% - each task is executed by all worker
	//   0% - no worker executes the tasks
	Spread  int
	RandRes RandomResult

	ws  []*taskengine.Worker
	wts taskengine.WorkerTasks
}

var rndPrice = RandomResult{
	Mean:   100,
	StdDev: 50,
}

// ========================================================

type Task struct {
	taskid string
	rndres *RandomResult
}

func (t *Task) TaskID() taskengine.TaskID { return taskengine.TaskID(t.taskid) }

// ========================================================

type demoResult struct {
	err      error
	price    float32
	currency string
}

func (res *demoResult) Error() error { return res.err }

func (res *demoResult) String() string {
	if res.err != nil {
		return "n/a"
	}
	return fmt.Sprintf("%.2f %s", res.price, res.currency)
}

// ========================================================

func demoWorkFn(ctx context.Context, worker *taskengine.Worker, workerInst int, task taskengine.Task) taskengine.Result {
	stask := task.(*Task)

	msec := stask.rndres.int64()

	res := &demoResult{}

	select {
	case <-ctx.Done():
		res.err = ctx.Err()
	case <-time.After(time.Duration(msec) * time.Millisecond):
		if !stask.rndres.success() {
			res.err = ErrorResult
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

func (scnr *Scenario) getWorkerID(n int) string {
	return getID("w", scnr.Workers, n)
}

func (scnr *Scenario) getTaskID(n int) string {
	return getID("t", scnr.Tasks, n)
}

func (scnr *Scenario) RandomWorkersAndTasks() error {

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
				task := &Task{
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
			task := &Task{
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

func (scnr *Scenario) ExecuteEvents() (chan *taskengine.Event, error) {

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
