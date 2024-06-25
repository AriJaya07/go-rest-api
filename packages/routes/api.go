package api

import (
	"log"
	"net/http"

	"github.com/AriJaya07/go-rest-api/packages/controllers/projects"
	"github.com/AriJaya07/go-rest-api/packages/controllers/tasks"
	"github.com/AriJaya07/go-rest-api/packages/controllers/users"
	"github.com/AriJaya07/go-rest-api/packages/models/store"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr  string
	store store.Store
}

func NewAPIServer(addr string, store store.Store) *APIServer {
	return &APIServer{
		addr:  addr,
		store: store,
	}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	projectService := projects.NewProjectService(s.store)
	projectService.RegisterRoutes(subrouter)

	usersService := users.NewUserService(s.store)
	usersService.RegisterRoutes(subrouter)

	tasksService := tasks.NewTasksService(s.store)
	tasksService.RegisterRoutes(subrouter)

	log.Println("Starting the API server at", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, subrouter))
}
