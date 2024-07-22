package storage

import "fmt"

func (s Storage) AddRecord(message, data string) (int64, error) {
	query := "INSERT INTO messages (message, data) VALUES ($1, $2) RETURNING id"
	var id int64
	err := s.DB.QueryRow(query, message, data).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AddToDb failed: %w", err)
	}
	return id, nil
}

func (s Storage) MessageConsumed(id int64) error {
	query := "UPDATE messages SET processed = true WHERE id = $1"
	_, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to mark message as consumed: %w", err)
	}
	return nil
}

func (s Storage) NumberOfProcessedMessages() (int, int, error) {
	var totalMessages, processedMessages int

	query := `
		SELECT
			(SELECT COUNT(*) FROM messages) AS total,
			(SELECT COUNT(*) FROM messages WHERE processed = true) AS processed
	`
	err := s.DB.QueryRow(query).Scan(&totalMessages, &processedMessages)
	if err != nil {
		return -1, -1, fmt.Errorf("error querying messages: %w", err)
	}

	return totalMessages, processedMessages, nil
}

func (s Storage) CloseDB() error {
	err := s.DB.Close()
	if err != nil {
		return err
	}
	return nil
}
