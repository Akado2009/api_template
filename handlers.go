package main

import (
	"api_template/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (env *Env) getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	userData, _, err := env.db.GetUser(userID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	returnJSON, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%s", returnJSON)
}

func (env *Env) saveUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	bodyString := string(body)

	var userData models.UserData
	err = json.Unmarshal([]byte(bodyString), &userData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	_, err = env.db.SaveUser(userData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%s", `{"result":"ok"}`)
}

func (env *Env) getAuthToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	bodyString := string(body)

	type AuthFormData struct {
		Password string `json:"password"`
		Login    string `json:"login"`
	}

	var afd AuthFormData
	err = json.Unmarshal([]byte(bodyString), &afd)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if afd.Login != "devers@inbox.ru" || afd.Password != "222" {
		fmt.Fprintf(w, "%s", `{"payload":"tdgdgYBatebe7453hsnY", "accepted":false}`)
	} else {
		fmt.Fprintf(w, "%s", `{"payload":"tdgdgarwtebe7453hsnY", "accepted":true}`)
	}

}

func (env *Env) getTestDataByToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "OPTIONS" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		return
	}

	authHeaderValue := r.Header.Get("Authorization")

	type AuthTokenData struct {
		Payload  string `json:"payload"`
		Accepted bool   `json:"accepted"`
	}

	var atd AuthTokenData
	err := json.Unmarshal([]byte(authHeaderValue), &atd)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Println(atd.Payload)
	fmt.Fprintf(w, "%s", `{"data":2}`)
}
