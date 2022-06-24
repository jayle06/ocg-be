package utils

import (
	"encoding/csv"
	"github.com/final-project/models"
	"log"
	"os"
	"strconv"
	"strings"
)

type ProductCSV struct {
	Name        string
	Description string
	Price       string
	IsSale      string
	PriceSale   string
	Images      []string
	Categories  []string
}

func ReadProductsCSV() []ProductCSV {
	csvFile, err := os.Open("import.csv")
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Println(err)
	}
	var productsCSV []ProductCSV
	for _, line := range csvLines {
		imgs := strings.Split(line[5], ";")
		cates := strings.Split(line[6], ";")
		product := ProductCSV{
			Name:        line[0],
			Description: line[1],
			Price:       line[2],
			IsSale:      line[3],
			PriceSale:   line[4],
			Images:      imgs,
			Categories:  cates,
		}
		productsCSV = append(productsCSV, product)
	}
	return productsCSV
}

func ConvertProductCSVToProduct(productsCSV []ProductCSV) []models.Product {
	var products []models.Product

	for i := 1; i < len(productsCSV); i++ {
		isSale, _ := strconv.ParseBool(productsCSV[i].IsSale)
		price, _ := strconv.ParseInt(productsCSV[i].Price, 0, 64)
		priceSale, _ := strconv.ParseInt(productsCSV[i].PriceSale, 0, 64)
		var images []models.Image
		var categories []models.Category
		for _, imgUrl := range productsCSV[i].Images {
			image := models.Image{
				ImageUrl: imgUrl,
			}
			images = append(images, image)
		}
		for _, cateId := range productsCSV[i].Categories {
			categoryID, _ := strconv.ParseInt(cateId, 0, 8)
			if categoryID != 0 {
				category := models.Category{
					ID: categoryID,
				}
				categories = append(categories, category)
			}
		}
		product := models.Product{
			Name:        productsCSV[i].Name,
			Description: productsCSV[i].Description,
			Price:       price,
			IsSale:      isSale,
			PriceSale:   priceSale,
			Images:      images,
			Categories:  categories,
		}
		products = append(products, product)
	}
	return products
}
