package timetracker_test

import (
	"fmt"
	"testing"
	"time"
	"timetracker"

	"github.com/google/go-cmp/cmp"
)

var tasklist = timetracker.Tasklist{

	"taskID1": {
		Name:        "piano",
		StartTime:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		ElapsedTime: time.Date(2021, 1, 1, 0, 10, 0, 0, time.UTC).Sub(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
	"taskID2": {
		Name:        "swim",
		StartTime:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		ElapsedTime: time.Date(2021, 1, 1, 0, 10, 0, 0, time.UTC).Sub(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestStartTracking(t *testing.T) {

	task := timetracker.NewTask("piano")

	if task.GetActive() {
		t.Error("task is already active")
	}

	task.Start(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))

	if !task.GetActive() {
		t.Error("task should be active")
	}

}

func TestStopTracking(t *testing.T) {

	task := timetracker.NewTask("piano")

	task.Start(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))

	if !task.GetActive() {
		t.Error("task should be active")
	}

	task.Stop(time.Date(2021, 1, 1, 0, 10, 0, 0, time.UTC))

	if task.GetActive() {
		t.Error("task should not  be active")
	}

}

func TestStartTime(t *testing.T) {

	task := timetracker.NewTask("piano")

	task.Start(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))

	if task.GetStartTime().IsZero() {
		t.Error("start time should not be zero")
	}

	task.Stop(time.Date(2021, 1, 1, 0, 10, 0, 0, time.UTC))

}

func TestElapsedTime(t *testing.T) {

	task := timetracker.NewTask("piano")

	task.Start(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))

	task.Stop(time.Date(2021, 1, 1, 0, 10, 0, 0, time.UTC))

	elapsed := task.GetElapsedTime()
	if elapsed == 0.0 {
		t.Error("elapsed time should not be zero")
	}

}

func TestGetMessage(t *testing.T) {

	name := "piano"

	task := timetracker.NewTask(name)

	task.Start(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))

	task.Stop(time.Date(2021, 1, 1, 0, 10, 0, 0, time.UTC))

	elapsed := task.GetElapsedTime()

	got := task.GetMessage()

	want := fmt.Sprintf("You spent %.1f seconds on the %s task", elapsed.Seconds(), name)

	if want != got {
		t.Errorf("Wanted: %s, got %s", want, got)
	}

}

func TestGetAllTasklist(t *testing.T) {
	want := []timetracker.Task{
		{
			Name: "piano",
		},
		{
			Name: "swim",
		},
	}

	got := tasklist.GetAllTaskList()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

// }

func TestSaveTask(t *testing.T) {

	tl := timetracker.Tasklist{}
	task := timetracker.NewTask("piano")
	task.Start(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))

	task.Stop(time.Date(2021, 1, 1, 0, 10, 0, 0, time.UTC))

	tl.Save(task)

	task2 := timetracker.NewTask("swim")

	task2.Stop(time.Date(2021, 1, 1, 0, 10, 0, 0, time.UTC))

	tl.Save(task2)

	got := tasklist.GetAllTaskList()
	want := tl.GetAllTaskList()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}
