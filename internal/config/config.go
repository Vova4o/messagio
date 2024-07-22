// Package config provides acces to get get env variables and flags for the program.
package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var flags = pflag.NewFlagSet("flags", pflag.ExitOnError)

func init() {
	// Define the flags and bind them to viper
	flags.StringP("ServerAddress", "a", "0.0.0.0:8080", "HTTP server network address")
	flags.StringP("Dsn", "p", "host=postgres port=5432 user=postgres password=password dbname=messages sslmode=disable timezone=UTC connect_timeout=5", "DB Connection string!")
	flags.StringP("Send", "s", "RequestTopic", "Reciver topic ID")
	flags.StringP("Recive", "r", "ResponseTopic", "Consumer topic ID")
	flags.StringP("Kafka", "k", "kafka:9092", "Kafka adress and port")
	flags.StringP("KafkaGroup", "g", "my-group", "Kafka group")

	// Parse the command-line flags
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Printf("Error parsing flags: %v", err)
	}

	// Bind the flags to viper
	bindFlagToViper("ServerAddress")
	bindFlagToViper("Dsn")
	bindFlagToViper("Send")
	bindFlagToViper("Recive")
	bindFlagToViper("Kafka")
	bindFlagToViper("KafkaGroup")

	// Set the environment variable names
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindEnvToViper("ServerAddress", "SERVICE_PORT")
	bindEnvToViper("Dsn", "DSN_STRING")
	bindEnvToViper("Send", "RECIVER_KAFKA_ID")
	bindEnvToViper("Recive", "CONSUMER_KAFKA_ID")
	bindEnvToViper("Kafka", "KAFKA_PORT")
	bindEnvToViper("KafkaGroup", "KAFKA_GROUP")

	// Read the environment variables
	viper.AutomaticEnv()
}

func bindFlagToViper(flagName string) {
	if err := viper.BindPFlag(flagName, flags.Lookup(flagName)); err != nil {
		log.Println(err)
	}
}

func bindEnvToViper(viperKey, envKey string) {
	if err := viper.BindEnv(viperKey, envKey); err != nil {
		log.Println(err)
	}
}

func Address() string {
	return viper.GetString("ServerAddress")
}

func Dsn() string {
	return viper.GetString("Dsn")
}

func Send() string {
	return viper.GetString("Send")
}

func Recive() string {
	return viper.GetString("Recive")
}

func Kafka() string {
	return viper.GetString("Kafka")
}

func KafkaGroup() string {
	return viper.GetString("KafkaGroup")
}
