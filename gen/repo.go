package gen

import (
	"database/sql"
	"fmt"
	"log"
	// _"github.com/go-sql-driver/mysql"
    _"github.com/lib/pq"
)

var REPO *Repo

type Repo struct{
	*Queries
	DB *sql.DB
}

type Cat struct{
	Legs int
	FurColor string
}

func LoadRepo(){
	var databaseUrl = "postgres://gakobo:Psql4321@localhost/ttnm_waste?sslmode=disable"
	fmt.Println("Connecting to database --- ",databaseUrl);
	connection,err := sql.Open("postgres",databaseUrl)
	if err!=nil{
		log.Fatal("Cannot connect to mysql database ",err)
	}
	REPO=&Repo{
		Queries: New(connection),
		DB: connection,
	}
}