package config

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestAddress(t *testing.T) {
	expected := ":8080"
	if address := Address(); address != expected {
		t.Errorf("expected Address to be %v, got %v", expected, address)
	}
}

func TestDsn(t *testing.T) {
	expected := "host=postgres port=5432 user=postgres password=password dbname=messages sslmode=disable timezone=UTC connect_timout=5"
	if dsn := Dsn(); dsn != expected {
		t.Errorf("expected DSN to be %v, got %v", expected, dsn)
	}
}

func TestMain(m *testing.M) {
	// Установка переменных окружения для тестирования
	os.Setenv("SERVICE_PORT", ":8080")
	os.Setenv("DSN_STRING", "host=postgres port=5432 user=postgres password=password dbname=messages sslmode=disable timezone=UTC connect_timout=5")

	// Инициализация viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	code := m.Run()

	// Очистка переменных окружения после тестирования
	os.Unsetenv("SERVICE_PORT")
	os.Unsetenv("DSN_STRING")

	os.Exit(code)
}
