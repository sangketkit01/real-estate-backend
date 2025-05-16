package main

import (
	"database/sql"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sangketkit01/real-estate-backend/api"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
	"github.com/sangketkit01/real-estate-backend/util"
)

const DB_SOURCE = "postgres://root:secret@localhost:5433/simple_real_estate?sslmode=disable"

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln(err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalln("cannot connect to DB:", err)
	}
	store := db.NewStore(conn)

	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatalln(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Server started at port 8080")
}
