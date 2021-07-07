package timetracker

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/template"
	"time"

	_ "github.com/lib/pq"
)

type Server struct {
	httpServer *http.Server
	Addr       string
	logger     *log.Logger
	Port       int
	LogLevel   string
	tasks      *Env
}

// type to hold options for Server struct
type Option func(*Server)

// override default port for Server
// provides cleaner user experience
func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func WithLogLevel(loglevel string) Option {
	return func(s *Server) {
		s.LogLevel = loglevel
	}
}

// server
func NewServer(opts ...Option) *Server {

	// create Server instance with defaults
	s := &Server{
		Port:     4000,
		LogLevel: "verbose",
	}

	// set override options.  loop takes in
	// With funcs loaded with input params and
	// executes to update Server struct
	for _, o := range opts {
		o(s)
	}

	newLogger := log.New(os.Stdout, "", log.LstdFlags)
	if s.LogLevel == "quiet" {
		newLogger.SetOutput(ioutil.Discard)
	}

	// update struct...perhaps there is a better way.
	s.Addr = fmt.Sprintf("127.0.0.1:%d", s.Port)
	s.logger = newLogger

	return s

}

func (s *Server) ListenAndServe() error {

	psqlInfo, err := BuildDbConnection()
	if err != nil {
		return fmt.Errorf("unable to build connection string: %s", err)
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	s.tasks = &Env{Db: db}

	s.httpServer = &http.Server{
		Addr:              s.Addr,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
		ErrorLog:          s.logger,
	}

	s.logger.Println("Starting up on ", s.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/task", s.showTaskReport)
	mux.HandleFunc("/task/create", s.createTask)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	s.httpServer.Handler = mux

	if err := s.httpServer.ListenAndServe(); err != nil {
		WaitForServerRoute(s.Addr + "/task")
		s.logger.Println("server start:", err)
		return err
	}

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c
	s.Shutdown()

	return nil
}

// waitForServerRoute checks if the main route is reachable
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

func (s *Server) Shutdown() {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	s.logger.Println("shutting down..")
	s.httpServer.Shutdown(ctx)
	os.Exit(0)

}

func GetEnvironmentVariable(env string) (string, error) {

	value := os.Getenv(env)

	if value == "" {
		return "", fmt.Errorf("%s value not set", env)
	}
	return value, nil
}

func BuildDbConnection() (string, error) {

	host, err := GetEnvironmentVariable("TIMETRACKER_DB_HOST")
	if err != nil {
		return "", fmt.Errorf("problem getting environment variable: %s", err)
	}
	port, err := GetEnvironmentVariable("TIMETRACKER_DB_PORT")
	if err != nil {
		return "", fmt.Errorf("problem getting environment variable: %s", err)
	}
	user, err := GetEnvironmentVariable("TIMETRACKER_DB_USER")
	if err != nil {
		return "", fmt.Errorf("problem getting environment variable: %s", err)
	}
	dbname, err := GetEnvironmentVariable("TIMETRACKER_DB_NAME")
	if err != nil {
		return "", fmt.Errorf("problem getting environment variable: %s", err)
	}

	convertPort, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("unable to convert to integer: %s", err)
	}

	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, convertPort, user, dbname), nil
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) showTaskReport(w http.ResponseWriter, r *http.Request) {

	report, err := s.tasks.GetReport()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%v", report)

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

	http.Redirect(w, r, "/task", http.StatusSeeOther)

}
