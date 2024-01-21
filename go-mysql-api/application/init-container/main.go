package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db_host string = os.Getenv("DB_HOST")
var db_user string = os.Getenv("DB_USER")
var db_pass string = os.Getenv("DB_PASS")
var db_port string = os.Getenv("DB_PORT")

type BytesBuffer struct {
	byteSlice []byte
}

func (b *BytesBuffer) Write(p []byte) (int, error) {
	b.byteSlice = p
	return len(b.byteSlice), nil
}

func main() {
	db := openDB()
	executeInitScript(db, readFile(os.Args[1]))

}

func openDB() *sql.DB {
	log.Println("Connecting to Datasource :", db_host+":"+db_port)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", db_user, db_pass, db_host, db_port)
	db, e := sql.Open("mysql", dsn)
	if e != nil {
		panic(e.Error())
	}
	log.Println("Connected successfully !")

	return db
}

func readFile(filename string) string {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0777)

	if err != nil {
		fmt.Println("An ERROR Occured while reading the file : ", err)
	}

	var b BytesBuffer

	mw := io.MultiWriter(&b)

	io.Copy(mw, f)

	return string(b.byteSlice)
}

func executeInitScript(db *sql.DB, query string) {
	queries := strings.Split(query, "\n")
	for i := range queries {
		if queries[i] == "" {
			continue
		} else {
			_, e := db.Query(queries[i])
			if e != nil {
				log.Fatal("An ERROR occured : ", e)
			}
		}

	}
	log.Println("Init script executed successfully !")
}
