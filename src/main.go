package main

import (
	"fmt"

	"github.com/rishabht08/pattern-mining/src/drain"
	"github.com/rishabht08/pattern-mining/src/kafkaconsumer"
)

func main() {
	logger := drain.New(drain.DefaultConfig())
	fmt.Println("Starting Kafka Patterns Consumer")
	kafkaconsumer.NewKafkaPatternsConsumer(logger)
}
