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

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tasks, err := s.tasks.GetLatest()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := TemplateData{Tasks: tasks}

	var ok bool

	data.PageTemplate, ok = s.templateCache["home.page.tmpl"]
	if !ok {
		fmt.Fprintf(w, fmt.Sprint("template does not exist: home.page.tmpl"))
		return
	}

	data.Render(w, r)
}

func (s *Server) showTaskReport(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/task/report" {
		http.NotFound(w, r)
		return
	}

	report, err := s.tasks.GetReport()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := TemplateData{Reports: report}

	var ok bool

	data.PageTemplate, ok = s.templateCache["report.page.tmpl"]
	if !ok {
		fmt.Fprintf(w, fmt.Sprintf("template does not exist: report.page.tmpl"))
		return
	}

	data.Render(w, r)

}

func (s *Server) createNewTaskForm(w http.ResponseWriter, r *http.Request) {

	data := TemplateData{}

	var ok bool

	data.PageTemplate, ok = s.templateCache["create.page.tmpl"]
	if !ok {
		fmt.Fprint(w, fmt.Sprintf("template does not exist: report.page.tmpl"))
		return
	}

	data.Render(w, r)

	fmt.Println("test")

	//http.Redirect(w, r, "/home", http.StatusSeeOther)

}

func (s *Server) startedTask(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, http.StatusBadRequest)
		return
	}

	newTask := r.Form.Get("task")

	task := NewTask(newTask)
	task.Start(time.Now())

	tasks := []Task{}
	tasks = append(tasks, task)

	data := TemplateData{Tasks: tasks}
	var ok bool

	data.PageTemplate, ok = s.templateCache["started.page.tmpl"]
	if !ok {
		fmt.Fprint(w, fmt.Sprintf("template does not exist: report.page.tmpl"))
		return
	}

	data.Render(w, r)

}

func (s *Server) stopTask(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, http.StatusBadRequest)
		return
	}

	newTask := r.Form.Get("task")
	log.Println(newTask)

	// task := NewTask(newTask)
	// task.Start(time.Now())

	// tasks := []Task{}
	// tasks = append(tasks, task)

	// data := TemplateData{Tasks: tasks}
	// var ok bool
	//fmt.Println(time.Parse(time.RFC3339, strStart))

	// startTime, err := time.Parse(time.RFC3339, strStart)
	// if err != nil {
	// 	fmt.Fprint(w, http.StatusBadRequest)
	// 	return
	// }

	// task := Task{
	// 	Name:      r.PostForm.Get("task"),
	// 	StartTime: startTime,
	// }
	// task.Stop(time.Now())
	// tasks := []Task{}

	// fmt.Println(task)
	// tasks = append(tasks, task)

	// data := TemplateData{Tasks: tasks}
	// var ok bool

	// data.PageTemplate, ok = s.templateCache["stop.page.tmpl"]
	// if !ok {
	// 	fmt.Fprint(w, fmt.Sprintf("template does not exist: report.page.tmpl"))
	// 	return
	// }

	// data.Render(w, r)

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
