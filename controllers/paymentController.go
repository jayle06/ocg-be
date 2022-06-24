package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/final-project/database"
	"github.com/final-project/utils"
	"net/http"
)

func CheckPaymentMethods(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	amount := data["total"]
	if data["payment"] == "MOMO" {
		pay := utils.MomoPayment(amount)
		json.NewEncoder(w).Encode(pay)
	}
}

func IPNMomo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("IPN received from Momo")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Method == "POST" {
		//json.NewEncoder(w).Encode(r.Form)
		fmt.Println(r.Form)
		var errCode = r.Form["errorCode"]
		var mess = r.Form["message"]
		var orderId = r.Form["orderId"]
		if errCode[0] == "0" {
			db := database.Connect()
			defer db.Close()
			_, err = db.Query("update orders set status_id = 2 where id = ?", orderId[0])
			if err != nil {
				return
			}
		}
		json.NewEncoder(w).Encode(map[string]string{
			"Order ID":   orderId[0],
			"Error Code": errCode[0],
			"Message":    mess[0],
		})
	}
}
