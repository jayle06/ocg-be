package repositories

import (
	"github.com/final-project/database"
	"github.com/final-project/models"
)

func CreateNewUser(user models.User) models.User {
	db := database.Connect()
	defer db.Close()
	_, err := db.Query("insert into users (name, phone_number, email, password, role_id) "+
		"values ( ?, ?, ?, ?, ?)", user.Name, user.PhoneNumber, user.Email, user.Password, 2)
	if err != nil {
		panic("Could not create new user")
	}
	return user
}

func GetAllUsers() []models.User {
	db := database.Connect()
	defer db.Close()
	rows, _ := db.Query("select u.id, u.name, u.phone_number, u.email, r.name " +
		"from users u join roles r on r.id = u.role_id")
	var users []models.User
	for rows.Next() {
		var user models.User
		rows.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.Role)
		users = append(users, user)
	}
	return users
}

func GetUserById(id int64) models.User {
	db := database.Connect()
	defer db.Close()
	var user models.User
	row, _ := db.Query("select u.id, u.name, u.phone_number, u.email, r.name "+
		"from users u join roles r on r.id = u.role_id where u.id = ? ", id)
	if row.Next() {
		row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.Role)
	}
	return user
}

func GetUserByEmail(email string) (models.User, error) {
	db := database.Connect()
	defer db.Close()
	var user models.User
	row, err := db.Query("select u.id, u.name, u.phone_number, u.email, r.name "+
		"from users u join roles r on r.id = u.role_id where u.email = ? ", email)
	if err != nil {
		return user, err
	}
	if row.Next() {
		row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.Role)
	}
	return user, nil
}

func UpdateUserById(id int64, user models.User) {
	db := database.Connect()
	defer db.Close()
	var roleId int64
	if user.Role == "ADMIN" {
		roleId = 1
	} else {
		roleId = 2
	}
	_, err := db.Query("update users set name = ?, phone_number = ?, email = ?, role_id = ? "+
		"where id = ?", user.Name, user.PhoneNumber, user.Email, roleId, id)
	if err != nil {
		return
	}
}

func DeleteUserById(id int64) {
	db := database.Connect()
	defer db.Close()
	_, err := db.Query("delete from users where id = ?", id)
	if err != nil {
		return
	}

}
