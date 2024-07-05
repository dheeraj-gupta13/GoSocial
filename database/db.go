package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "dheeraj"
	dbname   = "social"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	// defer db.Close()

	fmt.Println("Successfully connected!")

	// fmt.Println("SIRAJ", GetDB())
	if err = db.Ping(); err != nil {
		return nil, err
	}

	DB = db

	return DB, nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	fmt.Println("inside getbd*********", DB)
	return DB
}

// func initDB() {
// 	dsn := "your_mysql_user:your_mysql_password@tcp(127.0.0.1:3306)/your_database_name?charset=utf8&parseTime=True&loc=Local"
// 	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}
// 	db.AutoMigrate(&models.User{})
// }
