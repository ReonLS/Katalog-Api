package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	//open connection ke database
	driver, cString := GetConnectionString()
	db, err := sql.Open(driver, cString)

	if err != nil {
		fmt.Println("Database Gagal Connect: ", err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("DB enggak respon: ", err)
	}

	fmt.Println("Berhasil Konek")
	return db
}
