package main

import (
	"api_template/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (env *Env) sendRestorePasswordEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	//get base64 encrypted data
	vars := mux.Vars(r)
	linkData := vars["token"]

	tokenJSON, err := decryptTextAES256(linkData, env.crypto.AES256Key)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var token Token
	err = json.Unmarshal([]byte(tokenJSON), &token)
	if err != nil {
		fmt.Fprintf(w, "%s", getJSONAnswer("",
			false,
			"Некорректная ссылка!",
			""))
		return
	}

	if token.TTL-time.Now().Unix() > 0 {
		fmt.Fprintf(w, "%s", getJSONAnswer("",
			false,
			"Просроченная ссылка!",
			""))
		return
	}

	userData, _, err := env.db.GetUserByEmail(token.Email)
	if err != nil {
		fmt.Fprintf(w, "%s", getJSONAnswer("",
			false,
			"Некорректный email!",
			""))
		return
	}

	password, err := getNewPassword()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	userData.PswdHashB = getSHA256Bytes(password, env.crypto.SHA256Salt)

	userID, _, err := env.db.SaveUser(userData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if userID > 0 {

		err = sendEmail(env.smtp.Email, userData.Email, env.smtp.Host, env.smtp.Password,
			`Subject: Ваш пароль\n
			`+password)
		if err != nil {
			fmt.Fprintf(w, "%s", getJSONAnswer("",
				false,
				err.Error(),
				""))
			return
		}

		token, _ := encryptTextAES256Base64(getTokenJSON(userID, env.crypto.TokenTTL), env.crypto.AES256Key)

		fmt.Fprintf(w, "%s", getJSONAnswer(token,
			true,
			"",
			""))

	} else {
		fmt.Fprintf(w, "%s", getJSONAnswer("",
			false,
			"Ошибка смены пароля!",
			""))
	}

}

func (env *Env) sendRestorePasswordLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	var tokenItem Token
	err := converBody2JSON(r.Body, &tokenItem)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Println(tokenItem.Email)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	//check email format
	if !checkEmailFormat(tokenItem.Email) {
		fmt.Fprintf(w, "%s", getJSONAnswer("",
			false,
			"Некорректный EMail формат!",
			""))
		return
	}

	linkData, err := encryptTextAES256Base64(fmt.Sprintf(`{"email":"%s", "ttl":%d}`, tokenItem.Email, env.crypto.PasswordEmailTTL), env.crypto.AES256Key)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = sendEmail(env.smtp.Email, tokenItem.Email, env.smtp.Host, env.smtp.Password,
		`Subject: Смена пароля: `+env.crypto.RestorePasswordURL+linkData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%s", getJSONAnswer("",
		true,
		"Вам отправлен EMail с инструкциями!",
		""))
}

func (env *Env) registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	userData := new(models.UserData)
	err := converBody2JSON(r.Body, &userData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	userData.Email = strings.ToLower(userData.Email)

	//check email format
	if !checkEmailFormat(userData.Email) {
		fmt.Fprintf(w, "%s", getJSONAnswer("",
			false,
			"Некорректный EMail формат!",
			""))
		return
	}

	//generate password and send it to email
	password, err := getNewPassword()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	userData.PswdHashB = getSHA256Bytes(password, env.crypto.SHA256Salt)

	userID, errorCode, err := env.db.SaveUser(userData)
	if err != nil && errorCode != "22024" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if userID > 0 {

		err = sendEmail(env.smtp.Email, userData.Email, env.smtp.Host, env.smtp.Password,
			`Subject: Ваш пароль\n
			`+password)
		if err != nil {
			fmt.Fprintf(w, "%s", getJSONAnswer("",
				false,
				err.Error(),
				""))
			return
		}

		token, _ := encryptTextAES256Base64(getTokenJSON(userID, env.crypto.TokenTTL), env.crypto.AES256Key)

		fmt.Fprintf(w, "%s", getJSONAnswer(token,
			true,
			"",
			""))
	} else {
		fmt.Fprintf(w, "%s", getJSONAnswer("",
			false,
			"Такой EMail уже используется!",
			""))
	}
}

func (env *Env) getAuthToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	var afd AuthFormData
	err := converBody2JSON(r.Body, &afd)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	userData, _, err := env.db.GetUserByAuth(afd.Login, getSHA256Bytes(afd.Password, env.crypto.SHA256Salt))
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if userData.UserID > 0 {

		token, _ := encryptTextAES256Base64(getTokenJSON(userData.UserID, env.crypto.TokenTTL), env.crypto.AES256Key)

		fmt.Fprintf(w, "%s", getJSONAnswer(token,
			true,
			"",
			""))

	} else {

		fmt.Fprintf(w, "%s", getJSONAnswer("",
			false,
			"Не верный логин или пароль!",
			""))
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

	checked, err := checkAuthToken(r.Header.Get("Authorization"), env.crypto.AES256Key)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(405), 405)
		return
	}

	if !checked {
		fmt.Fprintf(w, "%s", getJSONAnswer("",
			false,
			"Невалидный токен!",
			""))

		return
	}

	fmt.Fprintf(w, "%s", getJSONAnswer("",
		true,
		"",
		`{"param":"value"}`))
}
