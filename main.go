package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var db *sql.DB

type studentInfo struct {
	Sid    string `json:"sid,omitempty"`
	Name   string `json:"name,omitempty"`
	Course string `json:"course,omitempty"`
}

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:8889)/grand?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func getStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	defer db.Close()
	ss := []studentInfo{}
	s := studentInfo{}
	rows, err := db.Query("SELECT * FROM student")
	if err != nil {
		fmt.Fprint(w, ""+err.Error())

	} else {
		for rows.Next() {
			rows.Scan(&s.Sid, &s.Name, &s.Course)
			ss = append(ss, s)
		}
		json.NewEncoder(w).Encode(ss)
	}

}

func addStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	defer db.Close()
	s := studentInfo{}

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	s.Sid = "Hi"
	fmt.Fprint(w, s)
	result, err := db.Exec("INSERT INTO student (name, course) values (?, ?)", s.Name, s.Course)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		_, err := result.LastInsertId()
		if err != nil {
			json.NewEncoder(w).Encode("{error: Record not inserted}")
		} else {
			json.NewEncoder(w).Encode(s)
		}
	}
}

func updateSudents(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "PUT")
}

func deleteStudents(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "DELETE")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students", addStudents).Methods("POST")
	r.HandleFunc("/students/{sid}", updateSudents).Methods("PUT")
	r.HandleFunc("/students/{sid}", deleteStudents).Methods("DELETE")
	http.ListenAndServe(":8080", r)
}
