package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/smtp"
)

//SMTPServerData struct for emails
type SMTPServerData struct {
	Host     string `json:"host"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func sendEmail(from string, to string, host string, password string, msg string) error {

	auth := smtp.PlainAuth("", from, password, host)

	if err := smtp.SendMail(host+":25", auth, from, []string{to}, []byte(msg)); err != nil {
		return err
	}

	return nil
}

func getJSONAnswer(token string, accepted bool, reason string, data string) string {
	return fmt.Sprintf(`{"accepted":%t, "token":"%s", "reason":"%s", "data":"%s"}`, accepted, token, reason, data)
}

func converBody2JSON(data io.Reader, v interface{}) error {
	body, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(string(body)), v)
	if err != nil {
		return err
	}

	return nil
}

func getNewPassword() (string, error) {
	guidBytes := make([]byte, 16)
	_, err := rand.Read(guidBytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", guidBytes[0:4]), nil
}
