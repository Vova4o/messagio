// Package service its a service level of aplication where all business logic created.
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/vova4o/messagio/internal/config"
	"github.com/vova4o/messagio/internal/models"
)

type Service struct {
	St Storager
}

var listOfQtoConsume = []string{
	"ResponseTopic",
}

type Storager interface {
	CloseDB() error
	AddRecord(string, string) (int64, error)
	MessageConsumed(int64) error                  // need to work on
	NumberOfProcessedMessages() (int, int, error) // need to work on
	Produce(string, []byte) error
	Consume() (string, []byte, error)
}

func NewService(st Storager) *Service {
	return &Service{St: st}
}

func (s *Service) GiveMeStats() (int, int, error) {
	totalMessages, processedMessages, err := s.St.NumberOfProcessedMessages()
	if err != nil {
		return -1, -1, err
	}
	return totalMessages, processedMessages, nil
}

func (s *Service) HandleMessage(msg models.MessageJSON) error {
	// we need to check if message is valid and send to DB
	queue := config.Send()

	if msg.Data == "" {
		return errors.New("message and data fields must not be empty")
	}

	dataString, err := convertAnyToString(msg.Data)
	if err != nil {
		return fmt.Errorf("error converting data to string: %w", err)
	}

	// add record to DB
	id, err := s.St.AddRecord(queue, dataString)
	if err != nil {
		return fmt.Errorf("error adding record: %w", err)
	}

	message := models.KafkaMessage{
		Id:      id,
		Message: dataString,
	}

	messageKAFKA, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error encoding message to JSON: %v", err)
	}

	// send data to kafka
	err = s.St.Produce(queue, messageKAFKA)
	if err != nil {
		return fmt.Errorf("kafka service error: %w", err)
	}

	return nil
}

func convertAnyToString(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", value), nil
	}
}

func (s *Service) CloseDB() error {
	err := s.St.CloseDB()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ConsumeMessagess() {
	for {
		topic, message, err := s.St.Consume()
		if err != nil {
			log.Println(err)
			continue
		}
		if topic == "RequestTopic" {
			err = s.St.Produce(listOfQtoConsume[0], message)
			if err != nil {
				log.Println(err)
				continue
			}
		} else if topic == "ResponseTopic" {
			var messageJSON models.KafkaMessage
			err = json.Unmarshal(message, &messageJSON)
			if err != nil {
				log.Println("Error decoding message:", err)
				continue
			}

			err = s.St.MessageConsumed(messageJSON.Id)
			if err != nil {
				log.Println("Db error:", err)
			}
		}
	}
}

func (s *Service) StartConsumeMessages() {
	go s.ConsumeMessagess()
}
