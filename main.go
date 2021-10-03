package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/snkzt/light-control-API/handler"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Create a new router
	router := mux.NewRouter()
	// Atttach paths
	router.HandleFunc("/lights", handler.GetLights).Methods("GET")
	router.HandleFunc("/lights/create", handler.CreateLight).Methods("POST")
	router.HandleFunc("/lights/update/{id:[a-zA-Z0-9]*}", handler.UpdateLight).Methods("PUT")
	router.HandleFunc("/lights/delete/{id:[a-zA-Z0-9]*}", handler.DeleteLight).Methods("DELETE")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeOut:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
