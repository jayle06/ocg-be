package repositories

import (
	"database/sql"
	"fmt"
	"github.com/final-project/database"
	"github.com/final-project/models"
	"log"
	"time"
)

func CreateOrder(order models.Order) (models.Order, error) {
	db := database.Connect()
	defer db.Close()

	row, err := db.Query("SELECT id FROM payments WHERE name = ?", order.Payment)
	if err != nil {
		log.Println(err)
		return order, err
	}

	var id int
	if row.Next() {
		row.Scan(&id)
	}

	_, err = db.Query("INSERT INTO orders (full_name, phone_number, email, address, total, payment_id, created_at) "+
		"VALUES (?, ?, ?, ?, ?, ?, NOW())", order.FullName, order.PhoneNumber, order.Email, order.Address, order.Total, id)
	if err != nil {
		log.Println(err)
		return order, err
	}

	row, err = db.Query("SELECT MAX(id) FROM orders")
	if err != nil {
		log.Println(err)
		return order, err
	}
	if row.Next() {
		row.Scan(&id)
	}
	fmt.Println(id)
	order.ID = int64(id)
	for _, product := range order.OrderItems {
		_, err = db.Query("INSERT INTO order_items VALUES (?, ? ,?)", product.ProductId, id, product.Quantity)
		if err != nil {
			log.Println(err)
			return order, err
		}
	}
	order.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	return order, nil
}

func GetAllOrders(month int, page int) ([]models.Order, error) {
	var orders []models.Order

	db := database.Connect()
	defer db.Close()
	var rows *sql.Rows
	var err error
	if month == 0 {
		rows, err = db.Query("SELECT o.id, o.full_name, o.total, p.name, o.created_at, s.name "+
			"FROM orders o JOIN payments p ON o.payment_id = p.id "+
			"JOIN status s WHERE o.status_id = s.id "+
			"LIMIT 10 OFFSET ?", page)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = db.Query("SELECT o.id, o.full_name, o.total, p.name, o.created_at "+
			"FROM orders o JOIN payments p ON o.payment_id = p.id "+
			"WHERE MONTH(o.created_at) = ? LIMIT 10 OFFSET ?", month, page)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	for rows.Next() {
		var order models.Order
		_ = rows.Scan(&order.ID, &order.FullName, &order.Total, &order.Payment, &order.CreatedAt, &order.Status)
		orders = append(orders, order)
	}
	return orders, nil
}

func GetRevenue(month int) (int, error) {
	db := database.Connect()
	defer db.Close()
	row, err := db.Query("SELECT SUM(total) FROM orders WHERE MONTH(created_at) = ?", month)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var total int
	if row.Next() {
		row.Scan(&total)
	}
	return total, nil
}

func GetTotalOrderByMonth(month int) (int, error) {
	db := database.Connect()
	defer db.Close()
	row, err := db.Query("SELECT COUNT(*) FROM orders WHERE MONTH(created_at) = ?", month)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var total int
	if row.Next() {
		row.Scan(&total)
	}
	return total, nil
}

func GetOrderById(id int) (models.Order, error) {
	db := database.Connect()
	defer db.Close()

	var order models.Order

	rows, err := db.Query("SELECT o.id, o.full_name, o.phone_number, o.email, o.address, o.total, p.name, o.created_at , status.id"+
		"FROM orders o JOIN payments p ON o.payment_id = p.id "+
		"JOIN status on status.id = o.status_id"+
		"WHERE o.id = ?", id)
	if err != nil {
		return order, err
	}
	if rows.Next() {
		rows.Scan(&order.ID, &order.FullName, &order.PhoneNumber, &order.Email, &order.Address, &order.Total, &order.Payment, &order.CreatedAt, &order.Status)
	}

	rows, err = db.Query("SELECT product_id, quantity FROM order_items WHERE order_id = ?", id)
	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		rows.Scan(&item.ProductId, &item.Quantity)
		items = append(items, item)
	}
	order.OrderItems = items
	return order, nil
}

func UpdateOrder(order models.Order) error {
	db := database.Connect()
	defer db.Close()
	_, err := db.Query("UPDATE orders "+
		"SET full_name = ?, email = ?, phone_number = ?, address = ?, total = ?, created_at = NOW() "+
		"WHERE id = ?", order.FullName, order.Email, order.PhoneNumber, order.Address, order.Total, order.ID)

	if err != nil {
		return err
	}
	_, err = db.Query("UPDATE orders "+
		"SET full_name = ?, email = ?, phone_number = ?, address = ?, total = ?, created_at = NOW() "+
		"WHERE id = ?", order.FullName, order.Email, order.PhoneNumber, order.Address, order.Total, order.ID)
	if err != nil {
		return err
	}
	_, _ = db.Query("DELETE FROM order_items WHERE id = ?", order.ID)

	for _, product := range order.OrderItems {
		_, err = db.Query("INSERT INTO order_items VALUES (?, ? ,?)", product.ProductId, order.ID, product.Quantity)
	}
	return nil
}

func DeleteOrder(id int) error {
	db := database.Connect()
	defer db.Close()
	_, err := db.Query("DELETE FROM order_items WHERE order_id = ?", id)
	if err != nil {
		return err
	}
	_, err = db.Query("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func GetCustomerOrders(email string, phone string, page int) ([]models.Order, error) {
	db := database.Connect()
	defer db.Close()
	var orders []models.Order
	rows, err := db.Query("SELECT o.id, o.full_name, o.email,o.total, p.name, o.created_at, status.name "+
		"FROM orders o JOIN payments p ON o.payment_id = p.id "+
		"JOIN status ON status.id = o.status_id "+
		"WHERE o.email = ? AND o.phone_number = ? LIMIT 10 OFFSET ?", email, phone, page)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.Order
		_ = rows.Scan(&order.ID, &order.FullName, &order.Email, &order.Total, &order.Payment, &order.CreatedAt, &order.Status)
		orders = append(orders, order)
	}
	return orders, nil
}

func GetRevenueByDay(offset int) (interface{}, error) {
	db := database.Connect()
	defer db.Close()

	//rows, err := db.Query()
	db.Query("DELETE FROM calendar")

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation).Format("2006-01-02")

	_, err := db.Query("CALL fill_calendar(?, DATE(NOW()))", firstOfMonth)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	rows, err := db.Query("select datefield, IFNULL(sum(total),0), IFNULL(sum(quantity), 0) " +
		"from calendar c " +
		"left join orders o on c.datefield = date(created_at) " +
		"left join order_items ot on o.id = ot.order_id  " +
		"where datefield  " +
		"group by c.datefield")
	if err != nil {
		log.Println(err)
		return 0, err
	}

	type Rev struct {
		Date          string `json:"date"`
		Revenue       int64  `json:"revenue"`
		TotalProducts int64  `json:"total_products"`
	}
	var revs []Rev

	for rows.Next() {
		var rev Rev
		rows.Scan(&rev.Date, &rev.Revenue, &rev.TotalProducts)
		revs = append(revs, rev)
	}
	return revs, nil
}
