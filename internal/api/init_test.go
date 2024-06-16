package api

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/ChillyWR/PasswordManager/config"
	"github.com/ChillyWR/PasswordManager/internal/repo"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	c, err := config.New()
	if err != nil {
		fmt.Println("Failed to init config")
	}

	testDB, err = repo.OpenConnection(&repo.Config{
		Host:     c.DB.Host,
		Port:     c.DB.Port,
		DBName:   c.DB.DBName,
		Username: c.DB.Username,
		SSLMode:  c.DB.SSLMode,
		Password: c.DB.Password,
	})
	if err != nil {
		fmt.Println("Failed to init db")
	}

	code := m.Run()

	db, err := testDB.DB()
	if err != nil {
		fmt.Println("Failed to get sql db")
	}
	db.Close()

	os.Exit(code)
}
