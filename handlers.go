package main

import (
	"api_template/models"
	"fmt"
	"net/http"
	"regexp"
)

func (env *Env) registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	var userData models.UserData
	err := converBody2JSON(r.Body, &userData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	//check email format
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(userData.Email) {
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

		token, _ := encryptTextAES256Base64(getTokenJSON(userID), env.crypto.AES256Key)

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

		token, _ := encryptTextAES256Base64(getTokenJSON(userData.UserID), env.crypto.AES256Key)

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
