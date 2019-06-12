package jobrunner

import (
	"time"

	"gopkg.in/robfig/cron.v2"
)

// StatusData encapsulates the relevant data regarding a single Job.
type StatusData struct {
	ID        cron.EntryID
	JobRunner *Job
	Next      time.Time
	Prev      time.Time
}

// Entries will return a detailed list of currently running recurring jobs
// to remove an entry, first retrieve the ID of entry
func Entries() []cron.Entry {
	return MainCron.Entries()
}

// StatusPage returns a slice of StatusData suitable for constructing a status
// page.
func StatusPage() []StatusData {

	ents := MainCron.Entries()

	Statuses := make([]StatusData, len(ents))
	for k, v := range ents {
		Statuses[k].ID = v.ID
		Statuses[k].JobRunner = AddJob(v.Job)
		Statuses[k].Next = v.Next
		Statuses[k].Prev = v.Prev

	}

	// t := template.New("status_page")

	// var data bytes.Buffer
	// t, _ = t.ParseFiles("views/Status.html")

	// t.ExecuteTemplate(&data, "status_page", Statuses())
	return Statuses
}

// StatusJSON returns a map with the slice of StatusData assigned to the
// key 'jobrunner'.
func StatusJSON() map[string]interface{} {

	return map[string]interface{}{
		"jobrunner": StatusPage(),
	}

}

// AddJob will cast a cron.Job to become a jobrunner.Job.
func AddJob(job cron.Job) *Job {
	return job.(*Job)
}
