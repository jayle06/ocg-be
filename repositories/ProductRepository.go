package repositories

import (
	"github.com/final-project/database"
	"github.com/final-project/models"
)

func GetAllProducts() []models.Product{
	db := database.Connect()
	defer db.Close()
	rows, _ := db.Query("SELECT * " +
		"FROM products")
	var products []models.Product
	for rows.Next(){
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
	rows, _ := db.Query("SELECT image_url FROM images " +
		"WHERE product_id = ?", id)
	var images []models.Image
	for rows.Next(){
		var image models.Image
		rows.Scan(&image.ImageUrl)
		images = append(images, image)
	}
	return images
}

func GetProductByID(id int64) models.Product{
	db := database.Connect()
	defer db.Close()
	rows, _ := db.Query("SELECT * FROM products " +
		"WHERE id = ?", id)
	var product models.Product
	if rows.Next(){
		rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.IsSale, &product.PriceSale, &product.CreatedAt)
	}
	product.Images = getImages(product.ID)
	return product
}

func UpdateProductByID( product models.Product) error {
	db := database.Connect()
	defer db.Close()
	_, err := db.Query("UPDATE products " +
		"SET name = ?, description = ?, price = ?, is_sale = ?, price_sale = ? " +
		"WHERE id = ?", product.Name, product.Description, product.Price, product.IsSale, product.PriceSale, product.ID)
	if err != nil{
		return err
	}
	return nil
}

func CreateProduct(product models.Product) error{
	db := database.Connect()
	defer db.Close()

	row, err := db.Query("SELECT MAX(id) FROM products")

	if err != nil{
		panic(err.Error())
	}
	if row.Next(){
		row.Scan(&product.ID)
	}else {
		product.ID = 0
	}
	product.ID += 1
	_, err = db.Query("INSERT INTO products " +
		"(name, description, price, is_sale, price_sale, created_at) " +
		"VALUES (?, ?, ?, ?, ?, NOW())", product.Name, product.Description, product.Price, product.IsSale, product.PriceSale)

	for _, image := range product.Images{
		//fmt.Println(image.ImageUrl)
		_, err = db.Query("INSERT INTO images VALUES (?, ?)", product.ID, image.ImageUrl)
	}
	if err != nil{
		return err
	}
	return nil
}

func DeleteProduct(id int64)error{
	db := database.Connect()
	defer db.Close()
	_, err := db.Query("DELETE FROM products " +
		"where id = ?", id)
	if err != nil{
		return err
	}
	return nil
}