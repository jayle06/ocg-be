package routers

import (
	"fmt"
	"github.com/final-project/controllers"
	"github.com/final-project/middlewares"
	"github.com/final-project/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func RunServer() {
	r := mux.NewRouter()
	r.Use(middlewares.CommonMiddleWare)

	customerURL := r.PathPrefix("/api/v1").Subrouter()
	customerURL.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	customerURL.HandleFunc("/products/{id}", controllers.GetProductByID).Methods("GET")
	customerURL.HandleFunc("/new-products", controllers.GetNewProducts).Methods("GET")
	customerURL.HandleFunc("/best-sale", controllers.GetBestSale).Methods("GET")
	customerURL.HandleFunc("/orders", controllers.CreateOrder).Methods("POST")
	customerURL.HandleFunc("/categories", controllers.GetAllCategories).Methods("GET")

	adminURL := r.PathPrefix("/admin").Subrouter()

	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("POST")

	adminURL.Use(middlewares.IsAuthenticated)

	adminURL.HandleFunc("/change-password", controllers.UpdatePassword).Methods("PUT")
	adminURL.HandleFunc("/update-information", controllers.UpdateInfo).Methods("PUT")
	adminURL.HandleFunc("/profiles", controllers.GetUserInfo).Methods("GET")

	adminURL.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	adminURL.HandleFunc("/products/{id}", controllers.GetProductByID).Methods("GET")
	adminURL.HandleFunc("/products/{id}", controllers.UpdateProductByID).Methods("PUT")
	adminURL.HandleFunc("/products", controllers.CreateProduct).Methods("POST")

	adminURL.HandleFunc("/orders", controllers.GetAllOrders).Methods("GET")
	adminURL.HandleFunc("/orders/{id}", controllers.GetOrderById).Methods("GET")
	adminURL.HandleFunc("/orders/{id}", controllers.UpdateOrder).Methods("PUT")
	adminURL.HandleFunc("/orders/{id}", controllers.DeleteOrder).Methods("DELETE")

	adminURL.HandleFunc("/categories", controllers.GetAllCategories).Methods("GET")
	adminURL.HandleFunc("/categories", controllers.CreateCategory).Methods("POST")
	adminURL.HandleFunc("/categories/{id}", controllers.UpdateCategory).Methods("PUT")
	adminURL.HandleFunc("/categories/{id}", controllers.DeleteCategory).Methods("DELETE")

	adminURL.HandleFunc("/uploads", utils.UploadFile).Methods("POST")
	images := http.StripPrefix("/images/", http.FileServer(http.Dir("./uploads/")))
	r.PathPrefix("/images/").Handler(images)

	admin := r.PathPrefix("/api/admin").Subrouter()

	admin.Use(middlewares.IsAuthorized)

	admin.HandleFunc("/products/{id}", controllers.DeleteProduct).Methods("DELETE")

	admin.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	admin.HandleFunc("/users", controllers.CreateNewUser).Methods("POST")
	admin.HandleFunc("/users/{id}", controllers.GetUserById).Methods("GET")
	admin.HandleFunc("/users/{id}", controllers.UpdateUserById).Methods("PUT")
	admin.HandleFunc("/users/{id}", controllers.DeleteUserById).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	//handler := cors.Default().Handler(r)
	fmt.Println("Server start on domain: http://localhost:10000")
	log.Fatal(http.ListenAndServe(":10000", handler))
}
