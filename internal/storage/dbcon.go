//Package storage creates connection to DB and Kafka
//Implements all necesary methods to work wit DB and Kafka.
package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/vova4o/messagio/internal/config"

	_ "github.com/jackc/pgx/v4/stdlib"
)

//Storage is a stucture where we hold DB connection and Kafka Producer and Consumer connections.
type Storage struct {
	DB       *sql.DB
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

var counts int64

//NewConn returns DB connection and Kafka connection.
func NewConn(brokers string, groupID string) *Storage {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres")
	}

	prod := createProducer(brokers)
	cons := createConsumer(brokers, groupID, getTopics())

	return &Storage{
		DB:       conn,
		Producer: prod,
		Consumer: cons,
	}
}

//createProducer subscribe to Kafka as producer
//brokers is connection string that comes from env var
func createProducer(brokers string) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		panic(err)
	}
	return p
}

//createConsumer subscribe to Kafka as consumer
//brokers is a connection string that comes from env var
//groupID name of the group that we connect to
//topics []string aray of topics that we want to subscribe
func createConsumer(brokers string, groupID string, topics []string) *kafka.Consumer {
    var c *kafka.Consumer
    var err error

    maxAttempts := 5
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        c, err = kafka.NewConsumer(&kafka.ConfigMap{
            "bootstrap.servers": brokers,
            "group.id":          groupID,
            "auto.offset.reset": "earliest",
        })

        if err == nil {
            err = c.SubscribeTopics(topics, nil)
            if err == nil {
                return c
            }
        }

        fmt.Printf("Attempt %d/%d failed: %v\n", attempt, maxAttempts, err)
        time.Sleep(time.Second * 5) // Подождите 5 секунд перед следующей попыткой
    }

    panic(fmt.Sprintf("Failed to create consumer after %d attempts: %v", maxAttempts, err))
}


//getTopics helpers function, provides lidt of topics that Kafka have to subscribe
func getTopics() []string {
	var top []string

	top = append(top, config.Send())
	top = append(top, config.Recive())

	return top
}

//openDB make connection to DB and returns it
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

//connectToDB helper function that check if DB connection created by trying to connect 
// to DB for up to 10 times, for slow servers can be adjusted to a higher number
func connectToDB() *sql.DB {
	dsn := config.Dsn()

	for {
		connection, err := openDB(dsn)
		if err != nil {
			fmt.Println(err)
			log.Println("Postgress not ready yet...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing of for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
