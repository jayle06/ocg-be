package controllers

import (
	"encoding/json"
	"github.com/final-project/models"
	repo "github.com/final-project/repositories"
	"github.com/final-project/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request){
	cat := r.FormValue("categoryId")
	orders := r.FormValue("_order")
	_page := r.FormValue("_page")
	page, _ := strconv.Atoi(_page)
	page = (page-1) * 12
	name := r.FormValue("q")
	name = "%" + name + "%"
	products, err := repo.GetAllProducts(cat, orders, page, name)
	if err != nil{
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "cannot get",
		})
		return
	}
	json.NewEncoder(w).Encode(products)
}

func GetNewProducts(w http.ResponseWriter, r *http.Request){
	utils.SetupResponse(&w, r)
	products := repo.GetNewProducts()
	json.NewEncoder(w).Encode(products)
}

func GetProductByID(w http.ResponseWriter, r *http.Request){
	utils.SetupResponse(&w, r)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	product := repo.GetProductByID(int64(id))
	json.NewEncoder(w).Encode(product)
}

func GetBestSale(w http.ResponseWriter, r *http.Request){
	utils.SetupResponse(&w, r)
	products := repo.GetBestSale()
	json.NewEncoder(w).Encode(products)
}

func UpdateProductByID(w http.ResponseWriter, r *http.Request){
	utils.SetupResponse(&w, r)
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	product.ID = int64(id)
	err = repo.UpdateProductByID(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product = repo.GetProductByID(product.ID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "success",
		"product updated": product,
	})
}

func CreateProduct(w http.ResponseWriter, r *http.Request){
	utils.SetupResponse(&w, r)
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println(product.Images)
	err = repo.CreateProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{ "message" : "success"})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request){
	utils.SetupResponse(&w, r)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	err := repo.DeleteProduct(int64(id))

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"message" : err.Error(),
		})
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message" : "success",
	})

}

