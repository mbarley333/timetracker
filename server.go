package timetracker

import (
	"context"
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
	"timetracker/ui"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type TaskStore interface {
	Create(task Task) (int, error)
	UpdateStopped(Task) error
	GetReport() ([]Report, error)
	GetLatest() ([]Task, error)
	GetTaskByName(string) (Task, error)
	GetTaskBySession() (Task, error)
	Delete(Task) error
	NewTaskSession(Task) error
}

type Server struct {
	httpServer    *http.Server
	Addr          string
	logger        *log.Logger
	Port          int
	LogLevel      string
	templateCache map[string]*template.Template
	TaskStore     TaskStore
}

// type to hold options for Server struct
type Option func(*Server) error

// override default port for Server
// provides cleaner user experience
func WithPort(port int) Option {
	return func(s *Server) error {
		s.Port = port
		return nil
	}
}

func WithNoLogging() Option {
	return func(s *Server) error {
		s.LogLevel = "quiet"
		return nil
	}
}

func WithPostgresStore(conn string) Option {
	return func(s *Server) error {

		db, err := NewPostgresStore(conn)
		if err != nil {
			return err
		}

		s.TaskStore = db
		return nil
	}
}

func WithSqliteStore() Option {
	return func(s *Server) error {

		db, err := NewSqliteStore()
		if err != nil {
			return err
		}

		s.TaskStore = db
		return nil
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
	s.Addr = fmt.Sprintf(":%d", s.Port)
	s.logger = newLogger

	return s

}

func (s *Server) ListenAndServe() error {

	var err error
	s.templateCache, err = NewTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	s.httpServer = &http.Server{
		Addr:              s.Addr,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
		ErrorLog:          s.logger,
	}

	s.httpServer.Handler = s.routes()
	s.logger.Println("Starting up on ", s.Addr)

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

func (s *Server) routes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.home)
	mux.HandleFunc("/task/report", s.showTaskReport)
	mux.HandleFunc("/task/create", s.createNewTaskForm)
	mux.HandleFunc("/task/started", s.startedTask)
	mux.HandleFunc("/task/stop", s.stopTask)

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("/static/", fileServer)

	return mux
}
