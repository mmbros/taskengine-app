package main

import (
	"fmt"
	"sort"

	"github.com/mmbros/taskengine"
)

// ================================================

type ByTaskID []taskengine.TaskID

func (s ByTaskID) Len() int {
	return len(s)
}

func (s ByTaskID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByTaskID) Less(i, j int) bool {
	return s[i] < s[j]
}

// ================================================

type ByWorkerID []taskengine.WorkerID

func (s ByWorkerID) Len() int {
	return len(s)
}

func (s ByWorkerID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByWorkerID) Less(i, j int) bool {
	return s[i] < s[j]
}

// ================================================

type TaskWorkers map[taskengine.TaskID][]taskengine.WorkerID

func GetTaskWorkers(wts taskengine.WorkerTasks) TaskWorkers {
	t2w := TaskWorkers{}

	for wid, ts := range wts {
		for _, t := range ts {
			tid := t.TaskID()
			wids, ok := t2w[tid]
			if ok {
				wids = append(wids, wid)
			} else {
				wids = []taskengine.WorkerID{wid}
			}
			t2w[tid] = wids
		}
	}

	// sort WorkerID list of each TaskID
	for tid, wids := range t2w {
		sort.Sort(ByWorkerID(wids))
		t2w[tid] = wids
	}

	return t2w
}

func (t2w TaskWorkers) Print() {
	// sort TaskID list
	tids := make([]taskengine.TaskID, 0, len(t2w))
	for key := range t2w {
		tids = append(tids, key)
	}
	sort.Sort(ByTaskID(tids))

	for k, tid := range tids {
		wids := t2w[tid]
		fmt.Printf("%2d) %s: %v\n", k+1, tid, wids)
	}

}
