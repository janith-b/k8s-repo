package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
)

var db_host string = os.Getenv("DB_HOST")
var db_user string = os.Getenv("DB_USER")
var db_pass string = os.Getenv("DB_PASS")
var db_port string = os.Getenv("DB_PORT")
var db_schema string = os.Getenv("DB_SCHEMA")
var db_table string = os.Getenv("DB_TABLE")

type UserInfo struct {
	UserName  string `json:"userName"`
	Dob       string `json:"dob"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func main() {
	db := openDB()
	http.HandleFunc("/addUser", insertData(db))
	http.HandleFunc("/getUsers", retrieveAllData(db))
	http.HandleFunc("/getUser", retrieveData(db))
	log.Fatal((http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))))
}

func openDB() *sql.DB {
	log.Println("Connecting to Datasource :", db_host+":"+db_port)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_user, db_pass, db_host, db_port, db_schema)
	db, e := sql.Open("mysql", dsn)
	if e != nil {
		panic(e.Error())
	}
	log.Println("Connected successfully !")

	return db
}

func insertDataMySQL(db *sql.DB, u UserInfo) error {
	query := fmt.Sprintf(`INSERT IGNORE INTO %s(userName,dob,firstName,lastName) VALUES("%s","%s","%s","%s");`, db_table, u.UserName, u.Dob, u.FirstName, u.LastName)
	_, e := db.Query(query)
	if e != nil {
		log.Fatal("An ERROR occured : ", e)
	}
	log.Println("Data inserted successfully")
	return e
}

func retrieveAllDataMySQL(db *sql.DB) []UserInfo {
	query := fmt.Sprintf(`SELECT * FROM %s;`, db_table)
	var allUsers []UserInfo
	var userInfo UserInfo
	rows, e := db.Query(query)
	if e != nil {
		log.Fatal("An ERROR occured : ", e)
	}
	for rows.Next() {
		rows.Scan(&userInfo.UserName, &userInfo.Dob, &userInfo.FirstName, &userInfo.LastName)
		allUsers = append(allUsers, userInfo)
	}
	return allUsers
}

func retrieveDataMySQL(db *sql.DB, userName string) (UserInfo, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE userName = "%s";`, db_table, userName)
	var userInfo UserInfo
	rows, e := db.Query(query)
	if e != nil {
		log.Fatal("An ERROR occured : ", e)
	}
	for rows.Next() {
		rows.Scan(&userInfo.UserName, &userInfo.Dob, &userInfo.FirstName, &userInfo.LastName)
	}

	return userInfo, e
}

func insertData(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var userInfo UserInfo
			bs, e := io.ReadAll(r.Body)
			if e != nil {
				log.Fatal("An ERROR occured : ", e)
			}
			json.Unmarshal(bs, &userInfo)
			if insertDataMySQL(db, userInfo) == nil {
				fmt.Fprintf(w, "User Inserted Successfuly")
			} else {
				fmt.Fprintf(w, "Failed Inserting User")
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)

		}
	}
}

func retrieveAllData(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

			bs, e := json.Marshal(retrieveAllDataMySQL(db))
			if e != nil {
				log.Println("An ERROR occured : ", e)
			}
			fmt.Fprintf(w, string(bs))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	}
}

func retrieveData(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userName := r.URL.Query().Get("userName")
			userInfo, e := retrieveDataMySQL(db, userName)
			if e != nil {
				log.Println("An ERROR occured : ", e)
			}
			bs, e := json.Marshal(userInfo)
			if e != nil {
				log.Println("An ERROR occured : ", e)
			}
			fmt.Fprintf(w, string(bs))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	}
}
