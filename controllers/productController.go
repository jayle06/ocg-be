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

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	cat := r.FormValue("categoryId")
	orders := r.FormValue("_order")
	_page := r.FormValue("_page")
	sale := r.FormValue("_sale")
	page, _ := strconv.Atoi(_page)
	page = (page - 1) * 12
	name := r.FormValue("q")
	name = "%" + name + "%"
	products, err := repo.GetAllProducts(cat, orders, page, name, sale)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	total, err := repo.GetTotalProductByRequest(cat, name, sale)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":    total,
		"products": products,
	})
}

func GetNewProducts(w http.ResponseWriter, r *http.Request) {
	products := repo.GetNewProducts()
	json.NewEncoder(w).Encode(products)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	product := repo.GetProductByID(int64(id))
	json.NewEncoder(w).Encode(product)
}

func GetBestSale(w http.ResponseWriter, r *http.Request) {
	products := repo.GetBestSale()
	json.NewEncoder(w).Encode(products)
}

func UpdateProductByID(w http.ResponseWriter, r *http.Request) {
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
		"message":         "success",
		"product updated": product,
	})
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = repo.CreateProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := repo.DeleteProduct(int64(id))

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func ImportProducts(w http.ResponseWriter, r *http.Request) {
	productsCSV := utils.ReadProductsCSV()
	products := utils.ConvertProductCSVToProduct(productsCSV)
	go func() {
		for _, product := range products {
			err := repo.CreateProduct(product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}()
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}
