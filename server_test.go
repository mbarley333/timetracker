package timetracker_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"timetracker"

	"github.com/google/go-cmp/cmp"
)

func TestRenderHomePage(t *testing.T) {
	t.Parallel()

	startTime, err := time.Parse(time.RFC3339, "2021-01-01T00:00:00+00:00")
	if err != nil {
		t.Fatal(err)
	}

	tasks := []timetracker.Task{
		{
			Name:           "piano",
			StartTime:      startTime,
			ElapsedTimeSec: 10.0,
		},
		{
			Name:           "swim",
			StartTime:      startTime,
			ElapsedTimeSec: 10.0,
		},
	}

	data := timetracker.TemplateData{Tasks: tasks}

	templateCache, err := timetracker.NewTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	var ok bool

	data.PageTemplate, ok = templateCache["home.page.tmpl"]
	if !ok {
		t.Error("template does not exist: home.page.tmpl")
	}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data.Render(w, r)
	}))

	WaitForServerRoute(ts.URL)

	client := ts.Client()

	rs, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		t.Fatal(err)
	}

	got, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/home_page_test.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
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

	data := timetracker.TemplateData{Reports: reports}

	templateCache, err := timetracker.NewTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	var ok bool

	data.PageTemplate, ok = templateCache["report.page.tmpl"]
	if !ok {
		t.Error("template does not exist: report.page.tmpl")

	}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data.Render(w, r)
	}))

	WaitForServerRoute(ts.URL)

	client := ts.Client()

	rs, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		t.Fatal(err)
	}

	got, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/report_page_test.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	if strings.TrimSpace(string(want)) != strings.TrimSpace(string(got)) {
		t.Errorf("want: %s, got: %s", string(want), string(got))
	}

}

func WaitForServerRoute(url string) {

	for {
		_, err := net.Dial("tcp", url)
		if err == nil {
			log.Println("tcp not listening")
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}

}
