package timetracker_test

import (
	"fmt"
	"testing"
	"time"
	"timetracker"
)

var startTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
var stopTime = time.Date(2021, 1, 1, 0, 10, 0, 0, time.UTC)

func TestStartTracking(t *testing.T) {

	task := timetracker.NewTask("piano")

	if task.GetActive() {
		t.Error("task is already active")
	}

	task.StartAt(startTime)

	if !task.GetActive() {
		t.Error("task should be active")
	}

}

func TestStopTracking(t *testing.T) {

	task := timetracker.NewTask("piano")

	task.StartAt(startTime)

	if !task.GetActive() {
		t.Error("task should be active")
	}

	task.Stop(stopTime)

	if task.GetActive() {
		t.Error("task should not  be active")
	}

}

func TestStartTime(t *testing.T) {

	task := timetracker.NewTask("piano")

	want := startTime

	task.StartAt(want)

	got := task.StartTime

	if !got.Equal(want) {
		t.Fatalf("want: %s, got:%s", want, got)
	}

}

func TestElapsedTime(t *testing.T) {

	task := timetracker.NewTask("piano")

	task.StartAt(startTime)

	task.Stop(stopTime)

	got := task.ElapsedTime.Seconds()
	want := 10 * time.Minute.Seconds()

	if want != got {
		t.Errorf("want: %f, got: %f", want, got)
	}

	got = task.ElapsedTimeSec

	if want != got {
		t.Errorf("want: %f, got: %f", want, got)
	}

}

func TestGetMessage(t *testing.T) {

	name := "piano"

	task := timetracker.NewTask(name)

	task.StartAt(startTime)

	task.Stop(stopTime)

	elapsed := task.ElapsedTime

	got := task.GetMessage()

	want := fmt.Sprintf("You spent %s seconds on the %s task", elapsed, name)

	if want != got {
		t.Errorf("Wanted: %s, got %s", want, got)
	}

}
