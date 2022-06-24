package middlewares

import (
	"encoding/json"
	"github.com/final-project/database"
	"github.com/final-project/models"
	"github.com/final-project/utils"
	"net/http"
	"strconv"
)

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		db := database.Connect()
		defer db.Close()

		row, _ := db.Query("select r.name from users u "+
			"join roles r on r.id = u.role_id where u.id = ?", user.ID)
		if row.Next() {
			row.Scan(&user.Role)
		}
		if user.Role != "ADMIN" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "access denied",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
