package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vova4o/messagio/internal/models"
	"github.com/vova4o/messagio/internal/service"
)

// MockStorager - мок интерфейса Storager для тестирования.
type MockStorager struct {
	mock.Mock
}

func (m *MockStorager) CloseDB() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorager) AddRecord(queue string, data string) (int64, error) {
	args := m.Called(queue, data)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorager) MessageConsumed(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStorager) NumberOfProcessedMessages() (int, int, error) {
	args := m.Called()
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockStorager) Produce(queue string, message []byte) error {
	args := m.Called(queue, message)
	return args.Error(0)
}

func (m *MockStorager) Consume() (string, []byte, error) {
	args := m.Called()
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}

// TestGiveMeStats проверяет корректность работы функции GiveMeStats.
func TestGiveMeStats(t *testing.T) {
	mockStorager := new(MockStorager)
	mockStorager.On("NumberOfProcessedMessages").Return(10, 5, nil)

	s := service.NewService(mockStorager)

	total, processed, err := s.GiveMeStats()

	assert.NoError(t, err)
	assert.Equal(t, 10, total)
	assert.Equal(t, 5, processed)
}

// TestHandleMessage проверяет обработку сообщений с валидными и невалидными данными.
func TestHandleMessage(t *testing.T) {
	mockStorager := new(MockStorager)
	mockStorager.On("AddRecord", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockStorager.On("Produce", mock.Anything, mock.Anything).Return(nil)

	s := service.NewService(mockStorager)

	t.Run("valid message", func(t *testing.T) {
		err := s.HandleMessage(models.MessageJSON{Data: "valid data"})
		assert.NoError(t, err)
	})

	t.Run("empty data", func(t *testing.T) {
		err := s.HandleMessage(models.MessageJSON{Data: ""})
		assert.Error(t, err)
		assert.Equal(t, "message and data fields must not be empty", err.Error())
	})
}
