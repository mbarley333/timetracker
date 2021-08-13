package timetracker_test

import (
	"testing"
	"time"
	"timetracker"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

func TestPostgres(t *testing.T) {

	// setup the machinery to run a test...in this case, we are using so why not spin it up
	// test the behavior to make sure the interface is working as expected
	t.Parallel()
	conn := "host=localhost port=5432 user=postgres dbname=timetracker sslmode=disable"

	var store timetracker.TaskStore
	var err error

	store, err = timetracker.NewPostgresStore(conn)
	if err != nil {
		t.Errorf("Error connecting to postgres: %s", conn)
	}

	taskname := "zzzzzzzz"
	_, err = store.GetTaskByName(taskname)
	if err != nil {
		t.Fatal(err)
	}

	task := timetracker.Task{
		Name: taskname,
	}

	var id int
	id, err = store.Create(task)
	if err != nil {
		t.Fatal(err)
	}

	task.Id = id

	err = store.NewTaskSession(task)
	if err != nil {
		t.Fatal(err)
	}

	got, err := store.GetTaskBySession()
	if err != nil {
		t.Fatal(err)
	}

	want := task

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	err = store.Delete(task)
	if err != nil {
		t.Fatal(err)
	}

}

func TestParseRowsTasks(t *testing.T) {

	t.Parallel()

	want := []timetracker.Task{
		{
			Name:           "piano",
			StartTime:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			ElapsedTimeSec: 10.0,
		},
		{
			Name:           "swim",
			StartTime:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			ElapsedTimeSec: 10.0,
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"task_name", "start_time", "elapsed_time"}).
		AddRow("piano", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), 10.0).
		AddRow("swim", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), 10.0)

	mock.ExpectQuery("SELECT task_name, start_time, elapsed_time FROM tasks ORDER BY start_time DESC LIMIT 10").WillReturnRows(rows)

	e := &timetracker.DBStore{Db: db}

	results, err := e.Db.Query(timetracker.LATEST)
	if err != nil {
		t.Fatal(err)
	}
	defer results.Close()

	got, err := timetracker.ParseRowsTasks(results)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestParseRowsReport(t *testing.T) {

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

	rows := sqlmock.NewRows([]string{"task_name", "total_time"}).
		AddRow("piano", 10).
		AddRow("swim", 10)

	mock.ExpectQuery("SELECT task_name, SUM(elapsed_time) total_time FROM tasks GROUP BY task_name ORDER BY SUM(elapsed_time) DESC").WillReturnRows(rows)

	e := &timetracker.DBStore{Db: db}

	results, err := e.Db.Query(timetracker.REPORT)
	if err != nil {
		t.Fatal(err)
	}
	defer results.Close()

	got, err := timetracker.ParseRowsReport(results)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}
