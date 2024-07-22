package kafkaconsumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	k2 "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rishabht08/pattern-mining/src/drain"
)

func NewKafkaPatternsConsumer(logger *drain.Drain) {
	servers := os.Getenv("BOOTSTRAP_SERVERS")
	if servers == "" {
		fmt.Println("Server values are unavailable")
		return
	}
	cfg := &k2.ConfigMap{
		"bootstrap.servers":               servers,
		"broker.address.family":           "v4",
		"group.id":                        "patterns_v2_logs",
		"partition.assignment.strategy":   "roundrobin",
		"auto.offset.reset":               "earliest",
		"go.application.rebalance.enable": true,
		"enable.auto.commit":              true,
		"session.timeout.ms":              10 * 1000,
	}

	client, err := k2.NewConsumer(cfg)
	if err != nil {
		fmt.Printf("Failed to create consumer: %v\n", err)
		return
	}

	defer client.Close()

	err = client.Subscribe("v3.db1.p2i13hg.log", nil)
	if err != nil {
		log.Fatal("Failed to start consumer", "err", err)
	}

	var wg sync.WaitGroup
	messages := make(chan *k2.Message, 100)
	done := make(chan bool)

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range messages {
				processAndProduce(msg, logger)
			}
		}()

	}

	// for {
	// 	msg, err := client.ReadMessage(10 * time.Second)
	// 	if err == nil {
	// 		if msg == nil || len(msg.Value) == 0 {
	// 			continue
	// 		}
	// 		messages <- msg
	// 	} else {
	// 		fmt.Printf("Consumer error: %v (%v)\n", err, msg)
	// 	}
	// }

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				msg, err := client.ReadMessage(10 * time.Second)
				if err == nil {
					messages <- msg
				} else {
					fmt.Printf("Consumer error: %v (%v)\n", err, msg)
					close(done)
					return
				}
			}
		}
	}()

	<-done
	close(messages)
	wg.Wait()

}

func processAndProduce(msg *k2.Message, logger *drain.Drain) {
	logs := bytes.Split(msg.Value, []byte("\n"))
	for _, logMsg := range logs {
		if len(logMsg) == 0 {
			continue
		}
		var logData []interface{}
		err := json.Unmarshal(logMsg, &logData)
		if err != nil {
			fmt.Printf("Invalid JSON message: %v\n", err)
			continue
		}
		if len(logData) < 8 {
			fmt.Println("Insufficient log data")
			continue
		}

		logCluster := logger.Train(logData[7].(string))
		fmt.Println((logCluster))

	}
}
