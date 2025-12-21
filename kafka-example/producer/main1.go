// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
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
// 	writer := &kafka.Writer{
// 		Addr:     kafka.TCP("kafka:9092"),
// 		Topic:    "orders",
// 		Balancer: &kafka.LeastBytes{},
// 	}
// 	defer writer.Close()

// 	ticker := time.NewTicker(2 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			order := Order{
// 				ID:      fmt.Sprintf("order-%d", rand.Intn(1000)),
// 				Product: fmt.Sprintf("item-%d", rand.Intn(10)),
// 				Qty:     rand.Intn(5) + 1,
// 				Time:    time.Now(),
// 			}

// 			data, _ := json.Marshal(order)
// 			msg := kafka.Message{
// 				Key:   []byte(order.ID),
// 				Value: data,
// 			}

// 			if err := writer.WriteMessages(context.Background(), msg); err != nil {
// 				log.Printf("Error: %v", err)
// 			} else {
// 				log.Printf("Produced: %+v", order)
// 			}
// 		}
// 	}
// }
