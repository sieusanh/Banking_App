package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/sieusanh/Banking_App/helpers"
)

// host, port, user, password, dbname
const (
	host = "localhost"
	port = "5433"
	user = "postgres"
	password = "1"
	dbname = "FintechBank"
)

// Variable DB is used to access DB in other files
var DB *gorm.DB

// Connect with DB, return DB connection object
func InitDatabase() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	database, err := gorm.Open("postgres", connStr)
	helpers.HandleErr(err)

	// Setup connection pool
	// Setup idle connections as 20
	database.DB().SetMaxIdleConns(20)
	// max connections as 200
	database.DB().SetMaxOpenConns(200)

	DB = database
}