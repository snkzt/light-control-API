package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DB struct {
	session    *mgo.Session
	collection *mgo.Collection
}

type Lights struct {
	ID     bson.ObjectId `json:"id" bson:"_id, omitempty"`
	Name   string        `json:"name" bson:"name, omitempty"`
	Status bool          `json:"status" bson:"status"`
}

func (db *DB) getLights(w http.ResponseWriter, r *http.Request) {
	var lights Lights
	c := db.session.DB("light").C("lights")
	err := c.Find(lights)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(lights)
		w.Write(response)
	}
}

func (db *DB) createLight(w http.ResponseWriter, r *http.Request) {
	var lights Lights
	postBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(postBody, &lights)
	// Create an Hash ID to insert
	lights.ID = bson.NewObjectId()
	c := db.session.DB("light").C("lights")
	err := c.Insert(lights)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(lights)
		w.Write(response)
	}
}

func (db *DB) updateLight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var lights Lights
	putBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(putBody, &lights)
	// Create an Hash ID to insert
	c := db.session.DB("light").C("lights")
	err := c.Update(bson.M{"_id": bson.ObjectIdHex(vars["id"])}, bson.M{"$set": &lights})
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text")
		w.Write([]byte("Updated successfully!"))
	}
}

func (db *DB) deleteLight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c := db.session.DB("light").C("lights")
	err := c.Remove(bson.M{"_id": bson.ObjectIdHex(vars["id"])})
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text")
		w.Write([]byte("Deleted successfully!"))
	}
}
