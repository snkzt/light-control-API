package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Light struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root@localhost:123asy123asy@tcp(127.0.0.1:3306)/lights_api")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/lights", getLights).Methods("GET")
	router.HandleFunc("/lights/create", createLights).Methods("POST")
	router.HandleFunc("/lights/update", updateLights).Methods("POST")
	router.HandleFunc("/lights/delete", deleteLights).Methods("POST")

	http.ListenAndServe(":8000", router)
}

func getLights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var lights []Light

	result, err := db.Query("SELECT id, name, status FROM lights")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var light Light
		err := result.Scan(&light.ID, &light.Name, &light.Status)
		if err != nil {
			panic(err.Error())
		}
		lights = append(lights, light)
	}

	json.NewEncoder(w).Encode(lights)
}

func createLights(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO lights(name, status) VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	status := keyVal["status"]

	_, err = stmt.Exec(name, status)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New light has created")
}

func updateLights(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE lights SET name = ?, status = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newName := keyVal["name"]
	newStatus := keyVal["status"]

	_, err = stmt.Exec(newName, newStatus, params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Light Name = %s has updated", params["name"])
}

func deleteLights(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM lights WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Light Name = %s has deleted", params["id"])
}
