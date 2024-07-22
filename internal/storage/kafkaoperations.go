package storage

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Produce отправляет сообщение в указанный топик Kafka.
// В случае ошибки при отправке сообщения возвращает ошибку.
//
// Аргументы:
//   topic - имя топика, в который будет отправлено сообщение.
//   msg - сообщение, которое необходимо отправить.
//
// Возвращаемое значение:
//   Возвращает ошибку, если отправка сообщения не удалась.
func (s *Storage) Produce(topic string, msg []byte) (error) {
	err := s.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
	if err != nil {
		return fmt.Errorf("Produce failed: %w", err)
	} else {
		return nil
	}
}

// Consume читает сообщение из Kafka.
// Блокирует выполнение до получения сообщения или возникновения ошибки.
//
// Возвращаемые значения:
//   topic - имя топика, из которого было прочитано сообщение.
//   message - прочитанное сообщение.
//   err - ошибка, возникшая во время чтения сообщения.
//
// В случае ошибки чтения сообщения возвращает ошибку.
func (s *Storage) Consume() (topic string, message []byte, err error) {
    for {
        msg, err := s.Consumer.ReadMessage(-1)
        if err == nil {
            return *msg.TopicPartition.Topic, msg.Value, nil
        } else {
            // The client will automatically try to recover from all errors.
            return "", nil, fmt.Errorf("consumer error: %v (%v)", err, msg)
        }
    }
}
