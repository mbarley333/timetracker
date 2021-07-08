package timetracker

import (
	"fmt"
	"time"
)

type Task struct {
	Name        string `db:"task_name"`
	Active      bool
	StartTime   time.Time `db:"start_time"`
	ElapsedTime float64   `db:"elasped_time"`
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

func (t Task) GetStartTime() time.Time {
	return t.StartTime
}

func (t *Task) Start(now time.Time) {
	t.Active = true
	t.StartTime = now
}

func (t *Task) Stop(now time.Time) {
	stop := now
	t.ElapsedTime = stop.Sub(t.StartTime).Seconds()
	t.Active = false

}

func (t Task) GetElapsedTime() float64 {
	return t.ElapsedTime
}

func (t Task) GetMessage() string {

	return fmt.Sprintf("You spent %.1f seconds on the %s task", t.ElapsedTime, t.Name)
}
