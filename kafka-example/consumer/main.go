package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	// ✅ SEPARATE consumers per topic
	coffeeWorker, err := ConnectConsumer([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer coffeeWorker.Close()

	teaWorker, err := ConnectConsumer([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer teaWorker.Close()

	fmt.Println("Consumers started (coffee_orders[0], tea_orders[0])")

	// ✅ Proper signal handling with context
	ctx, cancel := context.WithCancel(context.Background())
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigchan
		fmt.Println("Shutdown signal received")
		cancel()
	}()

	var wg sync.WaitGroup
	var msgCnt int32

	// helper to start consumers for all partitions of a topic
	startTopic := func(consumer sarama.Consumer, topic, label string) {
		partitions, err := consumer.Partitions(topic)
		if err != nil {
			fmt.Printf("%s: error fetching partitions: %v\n", label, err)
			return
		}
		if len(partitions) == 0 {
			fmt.Printf("%s: no partitions for topic %s\n", label, topic)
			return
		}
		for _, p := range partitions {
			pc, err := consumer.ConsumePartition(topic, p, sarama.OffsetOldest)
			if err != nil {
				fmt.Printf("%s: failed to consume partition %d: %v\n", label, p, err)
				continue
			}
			wg.Add(1)
			go func(pc sarama.PartitionConsumer, partition int32, lbl string) {
				defer wg.Done()
				defer pc.Close()
				for {
					select {
					case <-ctx.Done():
						fmt.Printf("%s consumer partition %d shutting down\n", lbl, partition)
						return
					case err := <-pc.Errors():
						fmt.Printf("%s error (partition %d): %v\n", lbl, partition, err)
					case msg := <-pc.Messages():
						cnt := atomic.AddInt32(&msgCnt, 1)
						fmt.Printf("%s #%d (topic=%s partition=%d): %s\n", lbl, cnt, msg.Topic, msg.Partition, string(msg.Value))
					}
				}
			}(pc, p, label)
		}
	}

	// start consumers for all partitions of each topic
	startTopic(coffeeWorker, "coffee_orders", "Coffee")
	startTopic(teaWorker, "tea_orders", "Tea")

	wg.Wait()
	fmt.Println("Processed", msgCnt, "messages")
}

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}
