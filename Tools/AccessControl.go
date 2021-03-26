package Tools

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type TokenHeader struct {
	Algorithm string `json:"alg"`
	TokenType string `json:"typ"`
}

type TokenData struct {
	Subject    string    `json:"sub"`
	Expiration time.Time `json:"exp"`
}

var DefaultUser = "Guest"

func Login(name, password string) (token string, err error) {
	user, err := DBGetUser(name)
	if err != nil {
		return "", err
	}
	if user.Password != password {
		return "", errors.New("Неверный пароль")
	}
	user.Active = true
	DBUpdateUser(user)
	return GetToken(user, time.Now().Add(168*time.Hour)), nil
}

func Logout(name string) {
	user, err := DBGetUser(name)
	if err != nil {
		fmt.Println("[Logout] ", err.Error())
		return
	}
	user.Active = false
	err = DBUpdateUser(user)
	if err != nil {
		fmt.Println("[Logout] ", err.Error())
	}
}

func GetToken(user DBUser, expiration time.Time) string {
	th, _ := json.Marshal(TokenHeader{"HS512", "JWT"})
	td, _ := json.Marshal(TokenData{user.Name, expiration})
	token := b64.StdEncoding.EncodeToString(th) + "." + b64.StdEncoding.EncodeToString(td)
	token += "." + b64.StdEncoding.EncodeToString([]byte(token))
	return token
}

func ParseToken(token string) (DBUser, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return DBUser{}, errors.New("Не правильный формат токена")
	}
	dec, err := b64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return DBUser{}, err
	}
	partsDec := strings.Split(string(dec), ".")
	if (len(partsDec) != 2) || (parts[0] != partsDec[0]) || (parts[1] != partsDec[1]) {
		return DBUser{}, errors.New("Не валидная подпись")
	}
	//th, err := b64.StdEncoding.DecodeString(parts[0])
	td, err := b64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return DBUser{}, err
	}
	var tds TokenData
	err = json.Unmarshal(td, &tds)
	if err != nil {
		return DBUser{}, err
	}
	if tds.Expiration.Before(time.Now()) {
		return DBUser{}, errors.New("Токен просрочен")
	}
	usr, err := DBGetUser(tds.Subject)
	if err != nil {
		return DBUser{}, errors.New("Токен не соответствует ни одному пользователю: " + tds.Subject)
	}
	if !usr.Active {
		fmt.Printf("[ParseToken] %v\n", usr)
		return DBUser{}, errors.New("Токен не актуален, выполнен выход")
	}
	return usr, nil
}
