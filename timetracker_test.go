package timetracker_test

import (
	"fmt"
	"testing"
	"time"
	"timetracker"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

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

	want := fmt.Sprintf("You spent %.1f seconds on the %s task", elapsed, name)

	if want != got {
		t.Errorf("Wanted: %s, got %s", want, got)
	}

}

// DB tests
func TestGenerateInsertSQL(t *testing.T) {
	t.Parallel()
	want := `INSERT INTO tasks(task_name, start_time, elased_time) VALUES($1, $2, $3)`

	got, err := timetracker.GenerateSQLQuery("insert")
	if err != nil {
		t.Fatal(err)
	}

	if want != got {
		t.Errorf("wanted: %s, got: %s", want, got)
	}

}

func TestGenerateReportSQL(t *testing.T) {
	t.Parallel()
	want := `SELECT task, SUM(elapsed_time) AS total_time FROM tasks GROUP BY task`

	got, err := timetracker.GenerateSQLQuery("report")
	if err != nil {
		t.Fatal(err)
	}

	if want != got {
		t.Errorf("wanted: %s, got: %s", want, got)
	}
}

func TestParseRows(t *testing.T) {

	t.Parallel()

	want := []timetracker.Report{
		{
			Task:      "piano",
			TotalTime: 10.0,
		},
		{
			Task:      "swim",
			TotalTime: 10.0,
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"task", "total_time"}).
		AddRow(10, "piano").
		AddRow(10, "swim")

	mock.ExpectQuery("SELECT task, SUM(elapsed_time) total_time FROM tasks GROUP BY task").WillReturnRows(rows)

	e := &timetracker.Env{Db: db}

	query, err := timetracker.GenerateSQLQuery("report")
	if err != nil {
		t.Fatal(err)
	}

	results, err := e.Db.Query(query)
	if err != nil {
		t.Fatal(err)
	}
	defer results.Close()

	got, err := timetracker.ParseRows(results)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}
