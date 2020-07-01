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
	db     models.Datastore
	crypto CryptoData
	smtp   SMTPServerData
}

func main() {

	connectionData := models.ConnectionData{}
	if gonfig.GetConf("config/db.json", &connectionData) != nil {
		log.Panic("load db confg error")
	}

	cryptoData := CryptoData{}
	if gonfig.GetConf("config/crypto.json", &cryptoData) != nil {
		log.Panic("load crypto confg error")
	}

	smtpServerData := SMTPServerData{}
	if gonfig.GetConf("config/smtp.json", &smtpServerData) != nil {
		log.Panic("load smtp confg error")
	}

	db, err := models.InitDB(connectionData.ToString())
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db, cryptoData, smtpServerData}

	router := mux.NewRouter()

	//login method
	router.HandleFunc("/getauthtoken/", env.getAuthToken).Methods("POST")
	//registration method
	router.HandleFunc("/registration/", env.registrationHandler).Methods("POST")
	//get an data with token example
	router.HandleFunc("/gettestdatabytoken/", env.getTestDataByToken).Methods("POST", "OPTIONS")

	//forgot password with sending special link

	//special rout for renew password and send it via email

	log.Fatal(http.ListenAndServe(":8080", router))
}
