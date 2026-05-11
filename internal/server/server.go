package server

import (
	"database/sql"
	"log"
	"net/http"

	"go_final_project/internal/handler"
	"go_final_project/internal/service"
	"go_final_project/internal/store"
)

func Init(port, webDir string, conn *sql.DB, logs *log.Logger) {
	taskStore := store.NewTaskStore(conn, logs)
	taskService := service.NewTaskService(taskStore, logs)
	taskHandler := handler.NewTaskHandler(taskService)

	utilsHandler := handler.NewUtilsHandler()

	srv := http.NewServeMux()
	srv.Handle("/", http.FileServer(http.Dir(webDir)))
	srv.Handle("/api/nextdate", utilsHandler.NextDate())
	srv.Handle("/api/tasks", taskHandler.List())
	srv.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.Get().ServeHTTP(w, r)
		case http.MethodPost:
			taskHandler.Add().ServeHTTP(w, r)
		case http.MethodPut:
			taskHandler.Edit().ServeHTTP(w, r)
		case http.MethodDelete:
			taskHandler.Delete().ServeHTTP(w, r)
		}
	})
	srv.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskHandler.Done().ServeHTTP(w, r)
		}
	})

	http.ListenAndServe(":"+port, srv)
}
