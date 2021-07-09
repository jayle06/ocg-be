package routers

import (
	"github.com/final-project/controllers"
	"github.com/final-project/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func RunServer() {
	r := mux.NewRouter()
	baseURL := r.PathPrefix("/api/v1").Subrouter()

	baseURL.HandleFunc("/login", controllers.Login).Methods("POST")
	baseURL.HandleFunc("/logout", controllers.Logout).Methods("POST")
	baseURL.HandleFunc("/change-password", controllers.UpdatePassword).Methods("PUT")
	baseURL.HandleFunc("/update-information", controllers.UpdateInfo).Methods("PUT")
	baseURL.HandleFunc("/profiles", controllers.GetUserInfo).Methods("GET")

	baseURL.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	baseURL.HandleFunc("/users", controllers.CreateNewUser).Methods("POST")
	baseURL.HandleFunc("/users/{id}", controllers.GetUserById).Methods("GET")
	baseURL.HandleFunc("/users/{id}", controllers.UpdateUserById).Methods("PUT")
	baseURL.HandleFunc("/users/{id}", controllers.DeleteUserById).Methods("DELETE")

	baseURL.HandleFunc("/uploads", utils.UploadFile).Methods("POST")
	images := http.StripPrefix("/images/", http.FileServer(http.Dir("./uploads/")))
	r.PathPrefix("/images/").Handler(images)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(":10000", handler))
}
