package jobrunner

import (
	"bytes"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"gopkg.in/robfig/cron.v2"
)

type Job struct {
	Name    string
	inner   cron.Job
	status  uint32
	Status  string
	running sync.Mutex
}

const UNNAMED = "(unnamed)"

func New(job cron.Job) *Job {
	name := reflect.TypeOf(job).Name()
	if name == "Func" {
		name = UNNAMED
	}
	return &Job{
		Name:  name,
		inner: job,
	}
}

func (j *Job) StatusUpdate() string {
	if atomic.LoadUint32(&j.status) > 0 {
		j.Status = "RUNNING"
		return j.Status
	}
	j.Status = "IDLE"
	return j.Status

}

func (j *Job) Run() {
	// If the job panics, just print a stack trace.
	// Don't let the whole process die.
	defer func() {
		if err := recover(); err != nil {
			//if revelError := revel.NewErrorFromPanic(err); revelError != nil {
			// revel.ERROR.Print(err, "\n", revelError.Stack)
			// } else {
			// revel.ERROR.Print(err, "\n", string(debug.Stack()))

			var buf bytes.Buffer
			logger := log.New(&buf, "JobRunner Log: ", log.Lshortfile)
			logger.Panic(err, "\n", string(debug.Stack()))
			// }
		}
	}()

	if !selfConcurrent {
		j.running.Lock()
		defer j.running.Unlock()
	}

	if workPermits != nil {
		workPermits <- struct{}{}
		defer func() { <-workPermits }()
	}

	atomic.StoreUint32(&j.status, 1)
	j.StatusUpdate()

	defer atomic.StoreUint32(&j.status, 0)
	defer j.StatusUpdate()

	j.inner.Run()
}