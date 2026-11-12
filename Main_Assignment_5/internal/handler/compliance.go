package handler

import (
	"aaa_bekend/internal/lmodel"
	"encoding/json"
	"net/http"
	"time"
)

func GetMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO 1: ambil profil user dari database
	profile, err := lmodel.GetProfileByUserID(user.ID)

	// TODO 2: handle error
	if err != nil {
		http.Error(w, "Gagal mengambil data profil", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func UpdateMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO 1: decode body JSON (username dan email)
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload tidak valid", http.StatusBadRequest)
		return
	}

	// TODO 2: update data user di database
	err := lmodel.UpdateUser(user.ID, req.Username, req.Email)
	if err != nil {
		http.Error(w, "Gagal memperbarui data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Data updated successfully",
		"username": req.Username,
	})
}

func GetMyAuditLog(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO 1: ambil audit log milik user
	logs, err := lmodel.GetAuditLogsByUserID(user.ID)

	// TODO 2: handle error
	if err != nil {
		http.Error(w, "Gagal mengambil audit log", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user.Username,
		"total": len(logs),
		"logs":  logs,
	})
}

func ExportMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO 1: ambil profil user
	profile, err := lmodel.GetProfileByUserID(user.ID)
	if err != nil {
		http.Error(w, "Gagal mengambil data profil untuk ekspor", http.StatusInternalServerError)
		return
	}

	// TODO 2: ambil project milik user
	projects, err := lmodel.GetProjectsByUserID(user.ID)
	if err != nil {
		http.Error(w, "Gagal mengambil data project untuk ekspor", http.StatusInternalServerError)
		return
	}

	// TODO 3: ambil audit log milik user
	logs, err := lmodel.GetAuditLogsByUserID(user.ID)
	if err != nil {
		http.Error(w, "Gagal mengambil audit log untuk ekspor", http.StatusInternalServerError)
		return
	}

	// TODO 4: gabungkan semua data jadi satu map ekspor
	export := map[string]interface{}{
		"exported_at": time.Now(),
		"profile":     profile,
		"projects":    projects,
		"audit_logs":  logs,
	}

	// TODO 5: set header download file dan kirim JSON
	w.Header().Set("Content-Disposition", "attachment; filename=my_data_export.json")
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(export); err != nil {
		http.Error(w, "Gagal mengenkode data ekspor", http.StatusInternalServerError)
		return
	}
}

func DeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO 1: hapus akun user beserta semua data terkait
	err := lmodel.DeleteUserAndRelatedData(user.ID)
	if err != nil {
		http.Error(w, "Gagal menghapus akun", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Your account and all data have been permanently deleted",
	})
}
