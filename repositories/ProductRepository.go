package repositories

import (
	"database/sql"
	"github.com/final-project/database"
	"github.com/final-project/models"
)

func GetAllProducts(category string, orders string, page int, name string) ([]models.Product, error) {
	db := database.Connect()
	defer db.Close()
	
	var rows *sql.Rows
	var err error
	var products []models.Product

	if orders == "desc" {
		if category == "" {
			rows, err = db.Query("SELECT * FROM product WHERE name LIKE ? ORDER BY price desc LIMIT ? OFFSET ?", name, 12, page)
		}else {
			rows, err = db.Query("SELECT p.* FROM products p JOIN product_category c ON p.id = c.product_id WHERE c.category_id = ? AND p.name LIKE ? ORDER BY p.price desc LIMIT ? OFFSET ?", category, name, 12, page)
		}
		
	}else {
		if len(category) == 0{
			rows, err = db.Query("SELECT * FROM product WHERE name LIKE ? ORDER BY price ASC LIMIT ? OFFSET ?",name, 12, page)
		}else {
			rows, err = db.Query("SELECT p.* FROM products p JOIN product_category c ON p.id = c.product_id  WHERE c.id = ? AND p.name LIKE ? ORDER BY p.price asc LIMIT ? OFFSET ?",category ,name, 12, page)
		}
	}

	if err != nil{
		return products, err
	}

	for rows.Next() {
		var product models.Product
		rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.IsSale, &product.PriceSale, &product.CreatedAt)
		//fmt.Println(product)
		product.Images = getImages(product.ID)
		products = append(products, product)
	}
	return products, nil
}

func GetNewProducts() []models.Product {
	db := database.Connect()
	defer db.Close()
	rows, _ := db.Query("SELECT * FROM products ORDER BY created_at DESC LIMIT 5")
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

func GetBestSale() []models.Product{
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
		//fmt.Println(image.ImageUrl)
		_, err = db.Query("INSERT INTO images VALUES (?, ?)", product.ID, image.ImageUrl)
	}
	if err != nil {
		return err
	}
	return nil
}

func DeleteProduct(id int64) error {
	db := database.Connect()
	defer db.Close()
	_, err := db.Query("DELETE FROM products "+
		"where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
