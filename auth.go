package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

//CryptoData struct
type CryptoData struct {
	AES256Key  string `json:"AES256Key"`
	SHA256Salt string `json:"SHA256Salt"`
}

//Token struct
type Token struct {
	UserID int   `json:"user_id"`
	TTL    int64 `json:"ttl"`
}

//AuthFormData struct for login/password sending
type AuthFormData struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

func checkAuthToken(authHeaderValue string, decrypetKey string) (bool, error) {

	tokenJSON, _ := decryptTextAES256(strings.ReplaceAll(authHeaderValue, `"`, ""), decrypetKey)

	var token Token
	err := json.Unmarshal([]byte(tokenJSON), &token)
	if err != nil {
		return false, err
	}

	return (token.TTL-time.Now().Unix() > 0), nil
}

func getTokenJSON(userID int) string {
	return fmt.Sprintf(`{"user_id":%d, "ttl":%d}`, userID, time.Now().Unix()+60)
}

func getSHA256Bytes(text string, salt string) []byte {
	h := sha256.New()
	h.Write([]byte(text + salt))
	return h.Sum(nil)
}

func encryptTextAES256Base64(textString string, keyString string) (string, error) {

	if len(keyString) != 32 {
		panic("too short key!")
	}
	text := []byte(textString)
	key := []byte(keyString)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	txt := gcm.Seal(nonce, nonce, text, nil)

	return b64.StdEncoding.EncodeToString([]byte(txt)), nil
}

func decryptTextAES256(encryptedBase64 string, keyString string) (string, error) {

	key := []byte(keyString)

	ciphertext, err := b64.StdEncoding.DecodeString(encryptedBase64) //[]byte(encryptedText)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", nil
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", nil
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", nil
	}

	return string(plaintext), nil
}
