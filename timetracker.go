package timetracker

import (
	"fmt"
	"time"
)

type Task struct {
	Id             int
	Name           string `db:"task_name"`
	Active         bool
	StartTime      time.Time `db:"start_time"`
	ElapsedTime    time.Duration
	ElapsedTimeSec float64 `db:"elapsed_time"`
}

type Report struct {
	Task      string
	TotalTime float64
}

func NewTask(task string) Task {
	t := Task{
		Name: task,
	}
	return t
}

func (t Task) GetActive() bool {
	return t.Active
}

func (t *Task) StartAt(now time.Time) {
	t.Active = true
	t.StartTime = now
}

func (t *Task) Stop(now time.Time) {
	stop := now
	t.ElapsedTime = stop.Sub(t.StartTime)
	t.ElapsedTimeSec = t.ElapsedTime.Seconds()
	t.Active = false
}

func (t Task) GetMessage() string {

	return fmt.Sprintf("You spent %s seconds on the %s task", t.ElapsedTime, t.Name)
}
