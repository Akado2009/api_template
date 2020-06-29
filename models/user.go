package models

import (
	"github.com/lib/pq"
)

// UserData structure
type UserData struct {
	UserID    int    `json:"user_id" db:"user_id"`
	IsActive  bool   `json:"is_active" db:"is_active"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	PswdHash  string `json:"pswd_hash" db:"pswd_hash"`
}

// GetUser method
func (db *DB) GetUser(userID int) (*UserData, pq.ErrorCode, error) {

	var errorCode pq.ErrorCode

	rows, err := db.Queryx("SELECT user_id, is_active, first_name, last_name, email, pswd_hash from users.user_get($1)", userID)
	defer rows.Close()
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			errorCode = err.Code
		}
		return nil, errorCode, err
	}

	userData := new(UserData)

	for rows.Next() {
		err = rows.StructScan(&userData)
		if err != nil {
			return nil, errorCode, err
		}
		break
	}

	return userData, errorCode, nil
}

// SaveUser method
func (db *DB) SaveUser(userData UserData) (pq.ErrorCode, error) {

	var errorCode pq.ErrorCode

	rows, err := db.Queryx("select * from users.user_save($1, $2, $3, $4, $5, $6)", userData.UserID, userData.IsActive, userData.FirstName, userData.LastName, userData.Email, userData.PswdHash)
	defer rows.Close()
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			errorCode = err.Code
		}
		return errorCode, err
	}

	return errorCode, nil
}
