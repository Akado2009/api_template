package main

import (
	"api_template/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
	"github.com/tkanos/gonfig"
)

type mockDB struct{}

func (mdb *mockDB) GetUser(userID int) (*models.UserData, pq.ErrorCode, error) {

	userData := new(models.UserData)
	userData.UserID = 1
	userData.IsActive = true
	userData.FirstName = "test"
	userData.LastName = "test"
	userData.Email = "test@test.test"

	var errorCode pq.ErrorCode

	return userData, errorCode, nil
}

func (mdb *mockDB) GetUserByAuth(email string, pswdHashB []byte) (*models.UserData, pq.ErrorCode, error) {

	userData := new(models.UserData)
	userData.UserID = 1
	userData.IsActive = true
	userData.FirstName = "test"
	userData.LastName = "test"
	userData.Email = "test@test.test"

	var errorCode pq.ErrorCode

	return userData, errorCode, nil
}

func (mdb *mockDB) SaveUser(userData models.UserData) (int, pq.ErrorCode, error) {

	var errorCode pq.ErrorCode

	//TODO make real save
	return 0, errorCode, nil
}

//TestGetTestDataByToken func
func TestGetTestDataByToken(t *testing.T) {

	cryptoData := CryptoData{}
	if gonfig.GetConf("config/crypto.json", &cryptoData) != nil {
		log.Panic("load crypto confg error")
	}

	smtpServerData := SMTPServerData{}
	if gonfig.GetConf("config/smtp.json", &smtpServerData) != nil {
		log.Panic("load smtp confg error")
	}

	env := Env{db: &mockDB{}, crypto: cryptoData, smtp: smtpServerData}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/gettestdatabytoken/", nil)

	token, _ := encryptTextAES256Base64(getTokenJSON(1), env.crypto.AES256Key)
	req.Header.Add("Authorization", token)

	http.HandlerFunc(env.getTestDataByToken).ServeHTTP(rec, req)

	expected := `{"accepted":true, "token":"", "reason":"", "data":"{"param":"value"}"}`
	if expected != rec.Body.String() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	}
}
