package main

import (
	"api_template/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"
)

//Env struct
type Env struct {
	db models.Datastore
}

func main() {

	connectionData := models.ConnectionData{}
	if gonfig.GetConf("config/db.json", &connectionData) != nil {
		log.Panic("load confg error")
	}

	db, err := models.InitDB(connectionData.ToString())
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db}

	router := mux.NewRouter()

	router.HandleFunc("/getauthtoken/", env.getAuthToken).Methods("POST")
	router.HandleFunc("/gettestdatabytoken/", env.getTestDataByToken).Methods("POST", "OPTIONS")

	router.HandleFunc("/saveuser/", env.saveUserHandler).Methods("POST")
	router.HandleFunc("/getuser/{user_id}", env.getUserHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
