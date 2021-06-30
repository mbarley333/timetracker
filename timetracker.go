package timetracker

import (
	"fmt"
	"time"
)

var (
	timeNow   = time.Now
	timeAfter = time.After
)

type Tasklist map[string]Task

type Task struct {
	Name        string
	Active      bool
	StartTime   time.Time
	ElapsedTime time.Duration
}

type Report struct {
	TaskName  string
	TotalTime time.Duration
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
	t.ElapsedTime = stop.Sub(t.StartTime)
	t.Active = false

}

func (tl *Tasklist) Save(task Task) {
	(*tl)[task.Name] = task
}

func (tl Tasklist) GetAllTaskList() []Task {
	ret := []Task{}

	for _, v := range tl {
		ret = append(ret, v)
	}

	return ret
}

func (t Task) GetElapsedTime() time.Duration {
	return t.ElapsedTime
}

func (t Task) GetMessage() string {
	return fmt.Sprintf("You spent %.1f seconds on the %s task", t.ElapsedTime.Seconds(), t.Name)
}
