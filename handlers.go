package timetracker

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
	"timetracker/ui"
)

// TemplateData is used to load struct
// data into the ui .tmpl files
type TemplateData struct {
	Reports      []Report
	Tasks        []Task
	PageTemplate *template.Template
}

func (a *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tasks, err := a.TaskStore.GetLatest()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := TemplateData{Tasks: tasks}

	var ok bool

	data.PageTemplate, ok = a.templateCache["home.page.tmpl"]
	if !ok {
		fmt.Fprintf(w, fmt.Sprint("template does not exist: home.page.tmpl"))
		return
	}

	data.Render(w, r)
}

func (a *Application) showTaskReport(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/task/report" {
		http.NotFound(w, r)
		return
	}

	report, err := a.TaskStore.GetReport()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := TemplateData{Reports: report}

	var ok bool

	data.PageTemplate, ok = a.templateCache["report.page.tmpl"]
	if !ok {
		fmt.Fprintf(w, fmt.Sprintf("template does not exist: report.page.tmpl"))
		return
	}

	data.Render(w, r)

}

func (a *Application) createNewTaskForm(w http.ResponseWriter, r *http.Request) {

	data := TemplateData{}

	var ok bool

	data.PageTemplate, ok = a.templateCache["create.page.tmpl"]
	if !ok {
		fmt.Fprint(w, fmt.Sprintf("template does not exist: report.page.tmpl"))
		return
	}

	data.Render(w, r)

}

func (a *Application) startedTask(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, http.StatusBadRequest)
		return
	}

	taskName := r.Form.Get("task")

	fmt.Printf("startedTask: %s", taskName)

	task := NewTask(taskName)
	task.StartAt(time.Now())

	tasks := []Task{}
	tasks = append(tasks, task)

	//set values in struct to persist data across HTML pages
	//a.taskid, err = a.tasks.Create(task)
	a.taskid, err = a.TaskStore.Create(task)
	if err != nil {
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}
	a.taskStartTime = task.StartTime

	data := TemplateData{Tasks: tasks}
	var ok bool

	data.PageTemplate, ok = a.templateCache["started.page.tmpl"]
	if !ok {
		fmt.Fprint(w, fmt.Sprintf("template does not exist: report.page.tmpl"))
		return
	}

	data.Render(w, r)

}

func (a *Application) stopTask(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, http.StatusBadRequest)
		return
	}

	taskName := r.PostForm.Get("task") //r.Form.Get("task")

	task := Task{
		Id:        a.taskid,
		Name:      taskName,
		StartTime: a.taskStartTime,
	}
	task.Stop(time.Now())

	err = a.TaskStore.UpdateStopped(task)
	if err != nil {
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	tasks := []Task{}
	tasks = append(tasks, task)

	data := TemplateData{Tasks: tasks}
	var ok bool

	data.PageTemplate, ok = a.templateCache["stop.page.tmpl"]
	if !ok {
		fmt.Fprint(w, fmt.Sprintf("template does not exist: report.page.tmpl"))
		return
	}

	data.Render(w, r)

}

func (td TemplateData) Render(w http.ResponseWriter, r *http.Request) {

	ts := td.PageTemplate

	err := ts.Execute(w, td)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/*.page.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFS(ui.Files, page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFS(ui.Files, "html/*.layout.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFS(ui.Files, "html/*.partial.tmpl")
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
