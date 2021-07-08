package routers

import (
	"github.com/final-project/controllers"
	"github.com/final-project/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer() {
	r := mux.NewRouter()
	baseURL := r.PathPrefix("/api/v1").Subrouter()
	baseURL.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	baseURL.HandleFunc("/users", controllers.CreateNewUser).Methods("POST")
	baseURL.HandleFunc("/users/{id}", controllers.GetUserById).Methods("GET")
	baseURL.HandleFunc("/users/{id}", controllers.UpdateUserById).Methods("PUT")
	baseURL.HandleFunc("/users/{id}", controllers.DeleteUserById).Methods("DELETE")

	baseURL.HandleFunc("/uploads", utils.UploadFile).Methods("POST")
	images := http.StripPrefix("/images/", http.FileServer(http.Dir("./uploads/")))
	r.PathPrefix("/images/").Handler(images)

	log.Fatal(http.ListenAndServe(":10000", r))
}
