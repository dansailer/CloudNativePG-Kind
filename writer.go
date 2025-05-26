package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type TickerData struct {
	Type      string `json:"type"`
	ProductID string `json:"product_id"`
	Time      string `json:"time"`
	Price     string `json:"price"`
	Volume    string `json:"volume_24h"`
}

const (
	user   = "appRoot"
	host   = "localhost"
	port   = "5432"
	dbname = "app"
	wsURL  = "wss://ws-feed.exchange.coinbase.com"
)

func main() {
	password := os.Getenv("APPDBROOT_PASSWORD")
	if password == "" {
		log.Fatal("Environment variable APPDBROOT_PASSWORD is not set")
	}
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to target database:", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ticker_data2 (
        id SERIAL PRIMARY KEY,
		time TIMESTAMP,
		product_id VARCHAR(50),
		price VARCHAR(50),
		VOLUME VARCHAR(50)
    )`)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	u, err := url.Parse(wsURL)
	if err != nil {
		log.Fatal("Error parsing WebSocket URL:", err)
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(u.String())
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer c.Close()

	subscribeMsg := map[string]interface{}{
		"type":        "subscribe",
		"product_ids": []string{"BTC-USD"},
		"channels":    []string{"ticker"},
	}
	err = c.WriteJSON(subscribeMsg)
	if err != nil {
		log.Fatal("Error sending subscription message:", err)
	}

	for {
		var msg TickerData
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}
		if msg.Type == "ticker" {
			_, err = db.Exec("INSERT INTO ticker_data2 (time, product_id, price, volume) VALUES ($1, $2, $3, $4)", msg.Time, msg.ProductID, msg.Price, msg.Volume)
			if err != nil {
				log.Println("Error inserting into database:", err)
				continue
			}
			fmt.Printf("Stored: Time=%s, Product=%s, Price=%s, Volue=%s\n", msg.Time, msg.ProductID, msg.Price, msg.Volume)
		}
	}
}
