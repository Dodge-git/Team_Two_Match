package config

import (
	"User/internal/models"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func SetUpDatabaseConnection() *gorm.DB {
	_ = godotenv.Load()
	
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbPass := os.Getenv("DB_PASS")

	dsn := fmt.Sprintf(
		"host=%v user=%v dbname=%v password=%v port=%v sslmode=disable",
		dbHost, dbUser, dbName, dbPass, dbPort,
	)

	  for i := 0; i < 20; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {

            sqlDB, err2 := db.DB()
            if err2 == nil {
                err2 = sqlDB.Ping()
                if err2 == nil {
                    fmt.Println("Database connected successfully")
					if err := db.AutoMigrate(&models.User{}); err != nil{
						panic("Failed to migrate database")
					}
                    return db
                }
            }
        }

        fmt.Println("Waiting for database...", err)
        time.Sleep(2 * time.Second)
    }
	if err != nil {
		panic("Can,t Connect to Database , contact tehnical support")
	}
	return db
}
