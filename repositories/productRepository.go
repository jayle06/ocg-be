package repositories

import (
	"database/sql"
	"github.com/final-project/database"
	"github.com/final-project/models"
)

func GetAllProducts(category string, orders string, page int, name string, sale string) ([]models.Product, error) {
	db := database.Connect()
	defer db.Close()

	var rows *sql.Rows
	var err error
	var products []models.Product

	if orders == "desc" {
		if category == "" {
			rows, err = db.Query("SELECT * FROM products "+
				"WHERE name LIKE ? "+
				"ORDER BY price desc LIMIT 12 OFFSET ?", name, page)
		} else {
			rows, err = db.Query("SELECT p.* FROM products p "+
				"JOIN product_category c ON p.id = c.product_id "+
				"WHERE c.category_id = ? AND p.name LIKE ? "+
				"ORDER BY p.price desc LIMIT 12 OFFSET ?", category, name, page)
		}
	} else {
		if len(category) == 0 {
			rows, err = db.Query("SELECT * FROM products "+
				"WHERE name LIKE ? "+
				"ORDER BY price asc LIMIT 12 OFFSET ?", name, page)
		} else {
			rows, err = db.Query("SELECT p.* FROM products p "+
				"JOIN product_category c ON p.id = c.product_id "+
				"WHERE c.category_id = ? AND p.name LIKE ? "+
				"ORDER BY p.price asc LIMIT 12 OFFSET ?", category, name, page)
		}
	}
	if err != nil {
		return products, err
	}
	for rows.Next() {
		var product models.Product
		rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price,
			&product.IsSale, &product.PriceSale, &product.CreatedAt)
		product.Images = getImages(product.ID)
		products = append(products, product)
	}
	return products, nil
}

func GetTotalProductByRequest(category string, name string, sale string) (int, error) {
	db := database.Connect()
	defer db.Close()

	var rows *sql.Rows
	var err error

	if category == "" {
		rows, err = db.Query("SELECT COUNT(*) FROM products "+
			"WHERE name LIKE ? ", name)
	} else {
		rows, err = db.Query("SELECT COUNT(*) FROM products p "+
			"JOIN product_category c ON p.id = c.product_id "+
			"WHERE c.category_id = ? AND p.name LIKE ? ", category, name)
	}

	if err != nil {
		return 0, err
	}
	var total int
	if rows.Next() {
		rows.Scan(&total)
	}
	return total, err
}

func GetNewProducts() []models.Product {
	db := database.Connect()
	defer db.Close()
	rows, _ := db.Query("SELECT * FROM products " +
		"ORDER BY created_at DESC LIMIT 5")
	var products []models.Product
	for rows.Next() {
		var product models.Product
		rows.Scan(&product.ID, &product.Name, &product.Description,
			&product.Price, &product.IsSale, &product.PriceSale, &product.CreatedAt)
		product.Images = getImages(product.ID)
		products = append(products, product)
	}
	return products
}

func getImages(id int64) []models.Image {
	db := database.Connect()
	defer db.Close()
	rows, _ := db.Query("SELECT image_url FROM images "+
		"WHERE product_id = ?", id)
	var images []models.Image
	for rows.Next() {
		var image models.Image
		rows.Scan(&image.ImageUrl)
		images = append(images, image)
	}
	return images
}

func GetProductByID(id int64) models.Product {
	db := database.Connect()
	defer db.Close()
	rows, _ := db.Query("SELECT * FROM products "+
		"WHERE id = ?", id)
	var product models.Product
	if rows.Next() {
		rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.IsSale, &product.PriceSale, &product.CreatedAt)
	}
	product.Images = getImages(product.ID)
	return product
}

func GetBestSale() []models.Product {
	db := database.Connect()
	defer db.Close()
	rows, _ := db.Query("SELECT p.* FROM products p LEFT JOIN order_items i on p.id = i.product_id GROUP BY p.id ORDER BY (sum(i.quantity)) desc LIMIT 5")
	var products []models.Product
	for rows.Next() {
		var product models.Product
		rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.IsSale, &product.PriceSale, &product.CreatedAt)
		//fmt.Println(product)
		product.Images = getImages(product.ID)
		products = append(products, product)
	}
	return products
}

func UpdateProductByID(product models.Product) error {
	db := database.Connect()
	defer db.Close()
	_, err := db.Query("UPDATE products "+
		"SET name = ?, description = ?, price = ?, is_sale = ?, price_sale = ? "+
		"WHERE id = ?", product.Name, product.Description, product.Price, product.IsSale, product.PriceSale, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func CreateProduct(product models.Product) error {
	db := database.Connect()
	defer db.Close()

	row, err := db.Query("SELECT MAX(id) FROM products")

	if err != nil {
		panic(err.Error())
	}
	if row.Next() {
		row.Scan(&product.ID)
	} else {
		product.ID = 0
	}
	product.ID += 1
	_, err = db.Query("INSERT INTO products "+
		"(name, description, price, is_sale, price_sale, created_at) "+
		"VALUES (?, ?, ?, ?, ?, NOW())", product.Name, product.Description, product.Price, product.IsSale, product.PriceSale)

	for _, image := range product.Images {
		_, err = db.Query("INSERT INTO images VALUES (?, ?)", product.ID, image.ImageUrl)
	}
	for _, category := range product.Categories {
		_, err = db.Query("INSERT INTO product_category VALUES (?, ?)", product.ID, category.ID)
	}
	if err != nil {
		return err
	}
	return nil
}

func DeleteProduct(id int64) error {
	db := database.Connect()
	defer db.Close()

	_, err := db.Query("DELETE FROM images "+
		"WHERE product_id = ?", id)
	if err != nil {
		return err
	}

	_, err = db.Query("DELETE FROM order_items "+
		"WHERE product_id = ?", id)
	if err != nil {
		return err
	}

	_, err = db.Query("DELETE FROM product_category "+
		"WHERE product_id = ?", id)
	if err != nil {
		return err
	}

	_, err = db.Query("DELETE FROM products "+
		"where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
