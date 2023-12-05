package internal

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/go-co-op/gocron"
)

type JobsState struct {
	Scheduler *gocron.Scheduler
	Peers     *PeerState
	DB        *sql.DB
}

// This is equivalent to the rows in the CronJobs table
type Job struct {
	ID         int
	Title      string
	Expression string
	Enabled    bool
}

// Adds a job
func (j JobsState) Add(title, expression string) error {
	result, err := j.DB.Exec("INSERT INTO CronJobs (title, expression, enabled) VALUES(?, ?, true)", title, expression)
	if err != nil {
		return fmt.Errorf("Could not insert: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Could not get insert ID: %s", err)
	}

	j.Schedule(int(id), expression)

	return nil
}

// Removes a job
func (j JobsState) Remove(id int) error {
	idStr := strconv.Itoa(id)

	err := j.Scheduler.RemoveByTag(idStr)
	if err != nil {
		return err
	}

	_, err = j.DB.Exec("DELETE FROM CronJobs WHERE id = ?;", idStr)
	return err
}

// Toggles a job
func (j JobsState) Toggle(id int) error {
	// Get the job based on the ID
	job := j.GetById(id)

	// Removes tag if the job is enabled
	if job.Enabled {
		err := j.Scheduler.RemoveByTag(strconv.Itoa(id))
		if err != nil {
			return err
		}
	} else {
		j.Schedule(id, job.Expression)
	}

	// Now set enabled to the opposi
	_, err := j.DB.Exec("UPDATE CronJobs SET enabled = NOT enabled WHERE id = ?;", id)
	return err
}

// Returns the jobs from the database
func (j JobsState) Get() []Job {
	rows, err := j.DB.Query("SELECT id, title, expression, enabled FROM CronJobs ORDER BY id;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var job_buf []Job
	for rows.Next() {
		var job Job
		job.Enabled = true
		if err := rows.Scan(&job.ID, &job.Title, &job.Expression, &job.Enabled); err != nil {
			log.Fatal(err)
		}
		job_buf = append(job_buf, job)
	}

	return job_buf
}

// Gets a job based on the id
func (j JobsState) GetById(id int) Job {
	var job Job

	jobs := j.Get()
	for _, jb := range jobs {
		if jb.ID == id {
			job = jb
		}
	}

	return job
}

// Re-adds all of the jobs to the scheduler
func (j JobsState) Persist() {
	jobs := j.Get()
	for _, job := range jobs {
		j.Schedule(job.ID, job.Expression)
	}
}

// Adds a cron job via the gocron scheduler
func (j JobsState) Schedule(id int, expression string) {
	j.Scheduler.Cron(expression).Tag(strconv.Itoa(id)).Do(func() {
		if err := Ring(id, j.Peers); err != nil {
			log.Fatal(err)
		}
	})
	j.Scheduler.StartAsync()
}
