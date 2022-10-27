package database

import "database/sql"

var db *sql.DB
var err error

func Connect() *sql.DB {
	db, err = sql.Open("mysql", "root:1234@tcp(localhost:3306)/secretdb?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Close() {
	if db != nil {
		db.Close()
	}
}
