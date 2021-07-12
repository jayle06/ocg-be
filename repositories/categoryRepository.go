package repositories

import (
	"github.com/final-project/database"
	"github.com/final-project/models"
)

func CreateCategory(category models.Category) error {
	db := database.Connect()
	defer db.Close()

	row, err := db.Query("SELECT MAX(id) FROM categories")
	if err != nil {
		panic(err.Error())
	}
	if row.Next() {
		row.Scan(&category.ID)
	} else {
		category.ID = 0
	}
	category.ID += 1

	_, err = db.Query("INSERT INTO categories VALUES (?, ?)", category.ID, category.Name)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCategory(category models.Category) error {
	db := database.Connect()
	defer db.Close()
	_, err := db.Query("DELETE FROM product_category "+
		"WHERE category_id = ?", category.ID)
	if err != nil {
		return err
	}
	_, err = db.Query("DELETE FROM categories "+
		"WHERE id = ?", category.ID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCategory(category models.Category) error {
	db := database.Connect()
	defer db.Close()

	_, err := db.Query("UPDATE categories SET name = ? WHERE id = ?", category.Name, category.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllCategories() ([]models.Category, error) {
	db := database.Connect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	var cates []models.Category
	for rows.Next() {
		var cate models.Category
		_ = rows.Scan(&cate.ID, &cate.Name)
		cates = append(cates, cate)
	}
	return cates, nil
}
