package controllers

import (
	"encoding/json"
	"github.com/final-project/database"
	"github.com/final-project/models"
	"github.com/final-project/utils"
	"net/http"
	"strconv"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := database.Connect()
	defer db.Close()

	var user models.User
	row, _ := db.Query("SELECT * FROM users "+
		"WHERE email = ?", data["email"])
	if row.Next() {
		row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.Password, &user.Role)
	}
	if user.ID == 0 {
		http.Error(w, "email incorrect", http.StatusBadRequest)
		return
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		http.Error(w, "password incorrect", http.StatusBadRequest)
		return
	}

	token, err := utils.GenerateJwt(strconv.Itoa(int(user.ID)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "login successfully",
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successfully",
	})
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, _ := utils.ParseJwt(cookie.Value)
	var user models.User
	db := database.Connect()
	defer db.Close()
	row, _ := db.Query("select u.id, u.name, u.phone_number, u.email, r.name "+
		"from users u join roles r on r.id = u.role_id where u.id = ? ", id)
	if row.Next() {
		row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.Role)
	}
	json.NewEncoder(w).Encode(user)
}

func UpdateInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, _ := utils.ParseJwt(cookie.Value)

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := database.Connect()
	defer db.Close()

	_, err = db.Query("update users set name = ?, phone_number = ?, email = ? "+
		"where id = ?", user.Name, user.PhoneNumber, user.Email, id)
	if err != nil {
		http.Error(w, "can not update information", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "update information successfully",
	})
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data["password"] != data["password_confirm"] {
		http.Error(w, "password confirm incorrect", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, _ := utils.ParseJwt(cookie.Value)
	userId, _ := strconv.Atoi(id)
	user := models.User{
		ID: int64(userId),
	}
	user.SetPassword(data["password"])

	db := database.Connect()
	defer db.Close()

	_, err = db.Query("update users set password = ? "+
		"where id = ?", user.Password, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "update password successfully",
	})
}
