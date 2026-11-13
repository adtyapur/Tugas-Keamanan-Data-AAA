package handler

import (
	"aaa_bekend/internal/db"
	"aaa_bekend/internal/lmodel"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type SubmitProjectReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GradeProjectReq struct {
	Score    int    `json:"score"`
	Feedback string `json:"feedback"`
}

// POST /api/v1/projects
func SubmitProject(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())

	// TODO 1: decode body JSON
	var req SubmitProjectReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO 2: validasi title dan description
	if req.Title == "" || req.Description == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	// TODO 3: simpan project ke database
	projectID, err := db.CreateProject(
		r.Context(),
		req.Title,
		req.Description,
		userCtx.ID,
	)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"project_id": projectID})
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	// TODO 1: ambil semua project dari database
	projects, err := db.GetAllProjects(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}

	// TODO 2: ganti nil jadi slice kosong
	if projects == nil {
		projects = []lmodel.Project{}
	}

	json.NewEncoder(w).Encode(&projects)
}

func GetAllProjectGrades(w http.ResponseWriter, r *http.Request) {
	// TODO 1: ambil semua nilai project dari database
	projects, err := db.GetAllProjectGrades(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch project grades", http.StatusInternalServerError)
		return
	}

	// TODO 2: ganti nil jadi slice kosong
	if projects == nil {
		projects = []lmodel.ProjectGrade{}
	}

	json.NewEncoder(w).Encode(&projects)
}

func GradeProject(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())

	// TODO 1: ambil project ID dari path
	projectID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// TODO 2: pastikan project ada
	_, err = db.GetProjectByID(r.Context(), projectID)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		return
	}

	// TODO 3: decode body JSON
	var req GradeProjectReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO 4: validasi score 0-100
	if req.Score < 0 || req.Score > 100 {
		http.Error(w, "Score must be between 0 and 100", http.StatusBadRequest)
		return
	}

	// TODO 5: simpan nilai ke database
	err = db.CreateGrade(
		r.Context(),
		projectID,
		userCtx.ID,
		req.Score,
		req.Feedback,
	)
	if err != nil {
		http.Error(w, "Failed to submit grade", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Grade Submitted"})
}

func GetProjectGrades(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())

	// TODO 1: ambil project ID dari path
	projectID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// TODO 2: ambil project, handle not found
	project, err := db.GetProjectByID(r.Context(), projectID)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		return
	}

	// TODO 3: student hanya boleh lihat nilai project miliknya sendiri
	if userCtx.Role == "student" {
		if project.StudentID != userCtx.ID {
			http.Error(w, "Forbidden: you can only view grades for your own projects", http.StatusForbidden)
			return
		}
	}

	// TODO 4: ambil nilai project dari database
	grades, err := db.GetGradesForProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Failed to fetch grades", http.StatusInternalServerError)
		return
	}

	// TODO 5: ganti nil jadi slice kosong
	if grades == nil {
		grades = []lmodel.Grade{}
	}
	json.NewEncoder(w).Encode(grades)
}
