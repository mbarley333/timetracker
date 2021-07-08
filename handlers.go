package timetracker

import (
	"log"
	"net/http"
	"text/template"
	"time"
)

// TemplateData is used to load struct
// data into the ui .tmpl files
type TemplateData struct {
	Reports []Report
	Tasks   []Task
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

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	Render(w, r, data, files)
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

	files := []string{
		"./ui/html/report.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	Render(w, r, data, files)

}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		return
	}

	data := Task{
		Name:        "test",
		StartTime:   time.Now(),
		ElapsedTime: 10.0,
	}

	err := s.tasks.Create(data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func Render(w http.ResponseWriter, r *http.Request, data TemplateData, files []string) {

	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
