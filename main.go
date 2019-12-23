package main

import (
	"./driver"
	dh "./handler/http"
	"./util"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
)

func main() {

	properties := util.NewProperties()

	dbName := properties.GetString("DB_NAME")
	dbHost := properties.GetString("DB_HOST")
	dbUser := properties.GetString("DB_USER")
	dbPort := properties.GetString("DB_PORT")

	connection, err := driver.ConnectSQL(dbHost, dbPort, dbUser, "", dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer connection.SQL.Close()

	initDataBase(connection.SQL)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	dHandler := dh.NewServerHandler(connection)
	r.Route("/", func(rt chi.Router) {
		rt.Mount("/scraping", domainRouter(dHandler))
	})

	http.ListenAndServe(":8005", r)

}

func domainRouter(dHandler *dh.Domain) http.Handler {
	r := chi.NewRouter()
	r.Get("/", dHandler.GetAllAddress)
	r.Get("/{address}", dHandler.GetByAddress)

	return r
}

func initDataBase(connection *sql.DB) {

	if _, err := connection.Exec(
		"CREATE TABLE IF NOT EXISTS domain (id SERIAL PRIMARY KEY, address varchar(100) NOT NULL, last_consultation TIMESTAMP NOT NULL)");
		err != nil {
		log.Fatal(err)
	}

	if _, err := connection.Exec(
		"CREATE TABLE IF NOT EXISTS detail_domain (id serial PRIMARY KEY, id_domain serial NOT NULL, ipaddress varchar(100) NOT NULL, servername varchar(200) NULL, grade varchar(10) NOT NULL, date TIMESTAMP NOT NULL, CONSTRAINT detail_domain_fk FOREIGN KEY (id_domain) REFERENCES domain(id))");
		err != nil {
		log.Fatal(err)
	}

	/*if _, err := connection.Exec(
		"INSERT INTO domain (address) VALUES ('google.com'), ('s4n.co')");
		err != nil {
		log.Fatal(err)
	}*/
}