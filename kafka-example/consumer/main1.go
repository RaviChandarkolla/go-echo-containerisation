// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/segmentio/kafka-go"
// )

// type Order struct {
// 	ID      string    `json:"id"`
// 	Product string    `json:"product"`
// 	Qty     int       `json:"qty"`
// 	Time    time.Time `json:"time"`
// }

// func main1() {
// 	reader := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:  []string{"kafka:9092"},
// 		Topic:    "orders",
// 		GroupID:  "order-processors",
// 		MinBytes: 10e3, // 10KB
// 		MaxBytes: 10e6, // 10MB
// 	})
// 	defer reader.Close()

// 	for {
// 		msg, err := reader.ReadMessage(context.Background())
// 		if err != nil {
// 			log.Printf("Error: %v", err)
// 			continue
// 		}

// 		var order Order
// 		json.Unmarshal(msg.Value, &order)

// 		fmt.Printf("Consumed: ID=%s, Product=%s, Qty=%d [Partition=%d, Offset=%d]\n",
// 			order.ID, order.Product, order.Qty, msg.Partition, msg.Offset)
// 	}
// }
