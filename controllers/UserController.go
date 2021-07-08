package controllers

import (
	"encoding/json"
	"github.com/final-project/models"
	repo "github.com/final-project/repositories"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if data["password"] != data["password_confirm"] {
		http.Error(w, "password confirm incorrect", http.StatusBadRequest)
		return
	}
	user := models.User{
		Name:        data["name"],
		PhoneNumber: data["phone_number"],
		Email:       data["email"],
	}
	user.SetPassword(data["password"])
	repo.CreateNewUser(user)

	response := make(map[string]string)
	response["message"] = "created"
	json.NewEncoder(w).Encode(response)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users := repo.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	user := repo.GetUserById(int64(id))
	json.NewEncoder(w).Encode(user)
}

func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	repo.UpdateUserById(int64(id), user)

	response := make(map[string]string)
	response["message"] = "updated"
	json.NewEncoder(w).Encode(response)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	repo.DeleteUserById(int64(id))
	response := make(map[string]string)
	response["message"] = "deleted"
	json.NewEncoder(w).Encode(response)
}
