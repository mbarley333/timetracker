package timetracker_test

import (
	"testing"
	"timetracker"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

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
	want := `SELECT task_name, SUM(elapsed_time) total_time FROM tasks GROUP BY task_name`

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

	mock.ExpectQuery("SELECT task_name, SUM(elapsed_time) total_time FROM tasks GROUP BY task_name").WillReturnRows(rows)

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
