package timetracker_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"timetracker"
)

func TestGetEnvironmentVariables(t *testing.T) {
	t.Parallel()

	type testCase struct {
		fn func(string) (string, error)
		a  string
	}
	tcs := []testCase{
		{fn: timetracker.GetEnvironmentVariable, a: "TIMETRACKER_DB_HOST"},
		{fn: timetracker.GetEnvironmentVariable, a: "TIMETRACKER_DB_PORT"},
		{fn: timetracker.GetEnvironmentVariable, a: "TIMETRACKER_DB_USER"},
		{fn: timetracker.GetEnvironmentVariable, a: "TIMETRACKER_DB_NAME"},
	}
	for _, tc := range tcs {
		_, err := tc.fn(tc.a)
		if err != nil {
			t.Fatalf("error with environment variable: %s %s", tc.a, err)
		}
	}

}

func TestServerRoutesStatusOK(t *testing.T) {
	t.Parallel()

	type testCase struct {
		route string
	}
	tcs := []testCase{
		{route: "http://127.0.0.1:9000"},
		{route: "http://127.0.0.1:9000/task/report"},
		{route: "http://127.0.0.1:9000/task/create"},
	}
	s := timetracker.NewServer(
		timetracker.WithPort(9000),
		timetracker.WithLogLevel("quiet"),
	)

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for _, tc := range tcs {
		resp, err := http.Get(tc.route)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unexpected status %d", resp.StatusCode)
		}

	}

}

func TestRenderHomePage(t *testing.T) {
	t.Parallel()

	tasks := []timetracker.Task{
		{
			Name:        "piano",
			StartTime:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			ElapsedTime: 10.0,
		},
		{
			Name:        "swim",
			StartTime:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			ElapsedTime: 10.0,
		},
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	data := timetracker.TemplateData{Tasks: tasks}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timetracker.Render(w, r, data, files)
	}))

	client := ts.Client()

	rs, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("piano")) {
		t.Errorf("want body to contain %q", []byte("piano"))
	}

	if !bytes.Contains(body, []byte("swim")) {
		t.Errorf("want body to contain %q", []byte("piano"))
	}

	if bytes.Contains(body, []byte("zzz")) {
		t.Errorf("do not want body to contain %q", []byte("zzz"))
	}

}

func TestRenderTaskReportPage(t *testing.T) {
	t.Parallel()

	reports := []timetracker.Report{
		{
			Task:      "piano",
			TotalTime: 10.0,
		},
		{
			Task:      "swim",
			TotalTime: 10.0,
		},
	}

	files := []string{
		"./ui/html/report.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	data := timetracker.TemplateData{Reports: reports}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timetracker.Render(w, r, data, files)
	}))

	client := ts.Client()

	rs, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("piano")) {
		t.Errorf("want body to contain %q", []byte("piano"))
	}

	if !bytes.Contains(body, []byte("swim")) {
		t.Errorf("want body to contain %q", []byte("piano"))
	}

	if bytes.Contains(body, []byte("zzz")) {
		t.Errorf("do not want body to contain %q", []byte("zzz"))
	}

}
