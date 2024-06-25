package projects

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/AriJaya07/go-rest-api/packages/controllers/auth"
	"github.com/AriJaya07/go-rest-api/packages/models/store"
	"github.com/AriJaya07/go-rest-api/packages/models/types"
	"github.com/AriJaya07/go-rest-api/packages/utils"
	"github.com/gorilla/mux"
)

type ProjectService struct {
	store store.Store
}

func NewProjectService(s store.Store) *ProjectService {
	return &ProjectService{store: s}
}

func (s *ProjectService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/projects", auth.WithJWTAuth(s.handleGetAllProject, s.store)).Methods("GET")
	r.HandleFunc("/projects/{id}", auth.WithJWTAuth(s.handleGetProject, s.store)).Methods("GET")
	r.HandleFunc("/projects", auth.WithJWTAuth(s.handleCreateProject, s.store)).Methods("POST")
	r.HandleFunc("/projects/{id}", auth.WithJWTAuth(s.handleDeleteProject, s.store)).Methods("DELETE")
}

func (s *ProjectService) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var project *types.Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if project.Name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Name is required"})
		return
	}

	err = s.store.CreateProject(project)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error creating project"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, project)
}

func (s *ProjectService) handleGetAllProject(w http.ResponseWriter, r *http.Request) {
	projects, err := s.store.GetAllProjects()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error getting all project"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, projects)
}

func (s *ProjectService) handleGetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	project, err := s.store.GetProject(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error getting project"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, project)
}

func (s *ProjectService) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := s.store.DeleteProject(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error deleting project"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
