package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	SQL * sql.DB
}


var dbConn = &DB{}


func ConnectSQL(host, port, user, pass, dbname string) (*DB, error){

	dbSource := fmt.Sprintf(
		"postgresql://%s@%s:%s/%s?sslmode=disable",
		user,
		host,
		port,
		dbname,
	)

	log.Println(dbSource)

	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
		panic(err)
	}

	dbConn.SQL = db
	return dbConn, err
}
