package internal

import (
	"database/sql"
	"errors"
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

	j.Scheduler.Cron(expression).Tag(strconv.Itoa(int(id))).Do(func() {
		if err := Ring(int(id), j.Peers); err != nil {
			log.Printf("Adding cron job failed: %s\n", err)
		}
	})
	j.Scheduler.StartAsync()

	return nil
}

// Removes a job
func (j JobsState) Remove(id string) error {
	err := j.Scheduler.RemoveByTag(id)
	_, err2 := j.DB.Exec("DELETE FROM CronJobs WHERE id = ?;", id)
	return errors.Join(err, err2)
}

// Removes a job
func (j JobsState) Toggle(id string) error {
	_ = j.Scheduler.RemoveByTag(id)
	_, err := j.DB.Exec("UPDATE CronJobs SET enabled = NOT enabled WHERE ID=?;", id)
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

// Re-adds all of the jobs to the scheduler
func (j JobsState) Persist() {
	jobs := j.Get()
	for _, job := range jobs {
		j.Scheduler.Cron(job.Expression).Tag(strconv.Itoa(job.ID)).Do(func() {
			if err := Ring(job.ID, j.Peers); err != nil {
				log.Fatal(err)
			}
		})
		j.Scheduler.StartAsync()
	}
}
