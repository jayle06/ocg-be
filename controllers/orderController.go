package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/final-project/models"
	repo "github.com/final-project/repositories"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	json.NewDecoder(r.Body).Decode(&order)
	x, err := repo.CreateOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(x)
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	_month := r.FormValue("_month")
	_page := r.FormValue("_page")
	var month, page int
	if _month == "" {
		month = 0
	} else {
		month, _ = strconv.Atoi(_month)
	}
	page, _ = strconv.Atoi(_page)
	page = (page - 1) * 10
	fmt.Println(month, page)
	orders, err := repo.GetAllOrders(month, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	totalRevenue, err := repo.GetRevenue(month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	total, err := repo.GetTotalOrderByMonth(month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_revenue": totalRevenue,
		"orders":        orders,
		"total":         total,
	})
}

func GetOrderById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	order, err := repo.GetOrderById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(order)
}
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := repo.DeleteOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var order models.Order
	json.NewDecoder(r.Body).Decode(&order)
	order.ID = int64(id)
	err := repo.UpdateOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func GetCustomerOrders(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("q")
	search = search + "%"
	_page := r.FormValue("_page")
	page, _ := strconv.Atoi(_page)
	page = (page - 1) * 10
	orders, err := repo.GetCustomerOrders(search, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(orders)
}
