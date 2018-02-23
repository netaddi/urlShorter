package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func getSqlConnection() *sql.DB {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/urlShorter")
	checkErr(err)
	return db
}

func directUrl(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	fmt.Println("access url: ", id)

	db := getSqlConnection()
	defer db.Close()

	rows, err := db.Query("SELECT url FROM urls WHERE id = ?", id)
	checkErr(err)

	var url string
	if rows.Next() {
		err = rows.Scan(&url)
		checkErr(err)
	}

	stmt, err := db.Prepare(
		"INSERT INTO access_log (url_id, time, remote_ip, forwarded_ip, UA, referer) VALUES (?, current_timestamp(), ?, ?, ?, ?)")
	checkErr(err)

	_, err = stmt.Exec(
		id,
		r.RemoteAddr,
		r.Header.Get("X-Forwarded-For"),
		r.Header.Get("User-Agent"),
		r.Header.Get("Referer"),
	)
	checkErr(err)

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

type Message struct {
	Url string `json:"url"`
}

func insertUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "POST" {

		// extract url to shorten
		body, err := ioutil.ReadAll(r.Body)
		checkErr(err)
		defer r.Body.Close()

		var msg Message
		err = json.Unmarshal(body, &msg)
		checkErr(err)
		url := msg.Url

		db := getSqlConnection()
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO urls (url) VALUES (?);")
		checkErr(err)

		res, err := stmt.Exec(url)
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)

		io.WriteString(w, strconv.Itoa(int(id)))
		fmt.Println("write url: ", id)

		stmt, err = db.Prepare(
			"INSERT INTO insert_log (url_id, time, remote_ip, forwarded_ip, UA, referer) VALUES (?, current_timestamp(), ?, ?, ?, ?)")
		checkErr(err)

		_, err = stmt.Exec(
			id,
			r.RemoteAddr,
			r.Header.Get("X-Forwarded-For"),
			r.Header.Get("User-Agent"),
			r.Header.Get("Referer"),
		)
		checkErr(err)
	}

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/u/{id}", directUrl)
	r.HandleFunc("/insert", insertUrl)
	http.Handle("/", r)

	fmt.Println(http.ListenAndServe(":8001", r))
}
