// Package queue contains nsq implementation
package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

const topic = "wallet"

// NSQ is a struct for queue producer.
type NSQ struct {
	producer *nsq.Producer
}

type message struct {
	Timestamp string  `json:"timestamp,omitempty"`
	Content   float64 `json:"content,omitempty"`
	Name      string  `json:"name,omitempty"`
}

// NewNSQ reads config and instantiates a producer.
func NewNSQ() (*NSQ, error) {
	config := nsq.NewConfig()

	log.Println("starting nsq producer at port 4150")
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		return nil, fmt.Errorf("NewNSQ error: %w", err)
	}

	err = producer.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot ping queue: %w", err)
	}

	return &NSQ{
		producer: producer,
	}, nil
}

// Stop gracefully stops the producer.
func (q *NSQ) Stop() {
	log.Println("gracefully stop the producer")
	q.producer.Stop()
}

// Wallet sends a message to the queue.
func (q *NSQ) Wallet(name string, ch chan error) {
	msg := message{
		Timestamp: time.Now().String(),
		Name:      name,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	err = q.producer.Publish(topic, payload)
	if err != nil {
		ch <- fmt.Errorf("cannot publish message to the queue: %w", err)
	}
	ch <- nil
}

// Operation sends a message to the queue.
func (q *NSQ) Operation(name string, content float64, ch chan error) {
	msg := message{
		Timestamp: time.Now().String(),
		Name:      name,
		Content:   content,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	err = q.producer.Publish(topic, payload)
	if err != nil {
		ch <- fmt.Errorf("cannot publish message to the queue: %w", err)
	}
	ch <- nil
}
