package main

import (
	"api_template/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
)

type mockDB struct{}

func (mdb *mockDB) GetUser(userID int) (*models.UserData, pq.ErrorCode, error) {

	userData := new(models.UserData)
	userData.UserID = 1
	userData.IsActive = true
	userData.FirstName = "test"
	userData.LastName = "test"
	userData.Email = "test@test.test"
	userData.PswdHash = "test"

	var errorCode pq.ErrorCode

	return userData, errorCode, nil
}

func (mdb *mockDB) SaveUser(userData models.UserData) (pq.ErrorCode, error) {
	//TODO make real save
	return "", nil
}

//TestGetQuestionsJSON func
func TestGetCatalogItemListJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)

	env := Env{db: &mockDB{}}
	http.HandlerFunc(env.getUserHandler).ServeHTTP(rec, req)

	expected := `{"user_id":1,"is_active":"true","first_name":"test","last_name":"test","email":"test@test.test","pswd_hash":"test"}`
	if expected != rec.Body.String() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	}
}
