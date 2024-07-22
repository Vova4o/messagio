package models

type MessageJSON struct {
	Data any `json:"data"`
} // @name MessageJSON

type KafkaMessage struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}

type GetStats struct {
	Total     int `json:"total"`
	Processed int `json:"processed"`
} // @name GetStats
