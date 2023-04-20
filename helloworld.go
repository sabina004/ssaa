package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	db *sql.DB
}

type OrderInfo struct {
	CustomerName   string
	CustomerEmail  string
	OrderTimestamp string
	TotalPrice     int
}

func dbConnect() server {
	db, err := sql.Open("sqlite3", "shop.db")
	if err != nil {
		log.Fatal(err)
	}

	s := server{db: db}

	return s
}

func (s *server) orderHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customerName := r.FormValue("customer_name")
	customerEmail := r.FormValue("customer_email")
	totalPrice := r.FormValue("total_price")

	totalPriceInt, err := strconv.Atoi(totalPrice)
	if err != nil {
		log.Fatal("totalPrice", err)
		return
	}

	orderTimestamp := time.Now().Format(time.RFC3339)

	_, err = s.db.Exec("INSERT INTO orders (customer_name, customer_email, order_date, total_price) VALUES (?, ?, ?, ?)", customerName, customerEmail, orderTimestamp, totalPriceInt)
	if err != nil {
		http.Error(w, "Failed to insert order", http.StatusInternalServerError)
		return
	}

	orderInfo := OrderInfo{
		CustomerName:   customerName,
		CustomerEmail:  customerEmail,
		OrderTimestamp: orderTimestamp,
		TotalPrice:     totalPriceInt,
	}

	outputHTML(w, "./static/postfinal.html", orderInfo)
}

func outputHTML(w http.ResponseWriter, filename string, orderInfo OrderInfo) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Fatal(err)
	}

	errExecute := t.Execute(w, orderInfo)
	if errExecute != nil {
		log.Fatal(errExecute)
	}
}

func main() {
	s := dbConnect()
	defer s.db.Close()

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	http.HandleFunc("/order", s.orderHandle)

	fmt.Println("Server running...")
	http.ListenAndServe(":8080", nil)
}
