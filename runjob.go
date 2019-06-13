// A job runner for executing scheduled or ad-hoc tasks asynchronously from HTTP requests.
//
// It adds a couple of features on top of the Robfig cron package:
// 1. Protection against job panics.  (They print to ERROR instead of take down the process)
// 2. (Optional) Limit on the number of jobs that may run simulatenously, to
//    limit resource consumption.
// 3. (Optional) Protection against multiple instances of a single job running
//    concurrently.  If one execution runs into the next, the next will be queued.
// 4. Cron expressions may be defined in app.conf and are reusable across jobs.
// 5. Job status reporting. [WIP]
package jobrunner

import (
	"time"

	"gopkg.in/robfig/cron.v2"
)

// Func is used to wrap a function literal as a job.
//
// For example:
//    jobrunner.Schedule("cron.frequent", jobs.Func(myFunc))
type Func func()

// Run will execute the underlying job logic
func (r Func) Run() { r() }

// Schedule will enqueue a job with the supplied name to execute according to
// the specification (spec) supplied.
func Schedule(spec string, name string, job cron.Job) error {
	sched, err := cron.Parse(spec)
	if err != nil {
		return err
	}
	MainCron.Schedule(sched, New(name, job))
	return nil
}

// Every will run the given job at a fixed interval. The interval provided is
// the time between the job ending and the job being run again.
// The time that the job takes to run is not included in the interval.
func Every(duration time.Duration, name string, job cron.Job) {

	MainCron.Schedule(cron.Every(duration), New(name, job))
}

// Now will run the given job immediately.
func Now(name string, job cron.Job) {
	go New(name, job).Run()
}

// In will run the given job once, after the given delay.
func In(duration time.Duration, name string, job cron.Job) {
	go func() {
		time.Sleep(duration)
		New(name, job).Run()
	}()
}
