package controllers

import (
	"encoding/json"
	"github.com/final-project/models"
	repo "github.com/final-project/repositories"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = repo.CreateCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := repo.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	category.ID = int64(id)
	err := repo.DeleteCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	category.ID = int64(id)
	err = repo.UpdateCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}
