package Tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

var tstList = []DBUser{{"AKai", "hitomi", 0, false}, {"guest", "", 10, false}, {"admin", "admin", 0, true}}
var tstNewUser = []User{{"AKai", "hitomi", 0, false}, {"guest", "", 10, false}, {"admin", "admin", 0, true}, {"newuser", "password", 10, false}}

func TestInit(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestInit")
	Init("../authdata.json")
	if !reflect.DeepEqual(userList, tstList) {
		t.Error("[TestInit] Загрузка пользователей не выполнена: ", userList)
	}
	fmt.Printf("[TestInit] Загруженные данные: %v\n[TestInit] Эталонный массив: %v\n", userList, tstList)
}

func TestGetUser(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestGetUser")
	userList = make([]User, 3)
	copy(userList, tstList)
	user, err := GetUser("admin")
	if (err != nil) || (user.Name != "admin") || (user.Password != "admin") || (user.Access != 0) {
		t.Error("[TestGetUser] Пользователь не найден: ", user)
	}
	fmt.Printf("[TestGetUser] Искомый пользователь: admin %v\n[TestFindUser] В массиве: %v\n", user, userList)
	user, err = GetUser("null")
	if err == nil {
		t.Error("[TestGetUser] Функция вернула данные пользователя null: ", user)
	}
	fmt.Printf("[TestGetUser] Искомый пользователь: null %v\n", user)
}

func TestNewUser(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestNewUser")
	userList = make([]User, 3)
	copy(userList, tstList)
	err := NewUser("newuser", "password", 10)
	if err != nil {
		t.Error("[TestNewUser] Ошибка добавления пользователя: ", err.Error())
	}
	if !reflect.DeepEqual(userList, tstNewUser) {
		t.Error("[TestNewUser] Пользователь не добавлен: ", userList)
	}
	fmt.Printf("[TestNewUser] Рабочий массив: %v\n[TestNewUser] Эталонный массив: %v\n", userList, tstNewUser)
	err = NewUser("newuser", "password", 10)
	if err == nil {
		t.Error("[TestNewUser] Не обработано добавление существующего пользователя")
	}
}

func TestRemoveUser(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestRemoveUser")
	userList = make([]User, 3)
	copy(userList, tstList)
	RemoveUser("guest")
	tst := []User{{"AKai", "hitomi", 0, false}, {"admin", "admin", 0, false}}
	if !reflect.DeepEqual(userList, tst) {
		t.Error("[TestRemoveUser] Пользователь не удален: ", userList)
	}
	fmt.Printf("[TestRemoveUser] Рабочий массив: %v\n[TestRemoveUser] Эталонный массив: %v\n", userList, tst)
}

func TestSaveUsers(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestSaveUsers")
	tst_filename := "tst_users.json"
	userList = make([]User, 3)
	copy(userList, tstList)
	if !SaveUsers(tst_filename) {
		t.Error("[TestSaveUsers] Функция вернула false")
	}
	txt, err := ioutil.ReadFile(tst_filename)
	if err != nil {
		t.Error("[TestSaveUsers] Не удалось прочитать файл: ", err.Error())
	}
	tst := make([]User, 0, 10)
	err = json.Unmarshal(txt, &tst)
	if err != nil {
		t.Error("[TestSaveUsers] Не удалось распарсить: ", err.Error())
	}
	if !reflect.DeepEqual(tstList, tst) {
		t.Error("[TestSaveUsers] Тест не пройден")
	}
	fmt.Printf("[TestSaveUsers] Сохраняемый массив: %v\n[TestSaveUsers] Прочитанный массив: %v\n", userList, tst)
}

func TestGetToken(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestGetToken")
	s := User{"admin", "admin", 0, false}
	token := GetToken(s, time.Date(2222, 1, 1, 0, 0, 0, 0, time.UTC))
	fmt.Println("[TestGetToken] ", token)
}

func TestParseToken(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestGetToken")
	userList = make([]User, 3)
	copy(userList, tstList)
	fmt.Println("[TestParseToken] Валидный токен")
	val, err := ParseToken("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhZG1pbiIsImV4cCI6IjIyMjItMDEtMDFUMDA6MDA6MDBaIn0=.ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnpkV0lpT2lKaFpHMXBiaUlzSW1WNGNDSTZJakl5TWpJdE1ERXRNREZVTURBNk1EQTZNREJhSW4wPQ==")
	tst := User{"admin", "admin", 0, true}
	if val != tst {
		fmt.Printf("[TestParseToken] Error: %v\n[TestParseToken] Получено: %v\n[TestParseToken] Ожидалось: %v\n", err, val, tst)
		t.Error("[TestParseToken] Валидный токен не расшифрован")
	}

	fmt.Println("[TestParseToken] Не валидный токен")
	val, err = ParseToken("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhZG1pbiIsImV4cCI6IjE5NzAtMDEtMDFUMDM6MDA6MDArMDM6MDAifQ==.ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnpkV0lpT2lKaFpHMXBiaUlzSW1WNGNDSTZJakU1TnpBdE1ERXRNREZVTURNNk1EQTZNREFyTURNNk1EQWlmUT09")

	if err == nil {
		fmt.Printf("[TestParseToken] Error: %v\n[TestParseToken] Session: %v\n", err, val)
		t.Error("[TestParseToken] Не валидный токен не вернул ошибки")
	}
}

func TestLogin(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestLogin")
	userList = make([]User, 3)
	copy(userList, tstList)
	token, err := Login("admin", "admin")
	if err != nil {
		t.Error("[TestLogin] ", err.Error())
	}
	fmt.Printf("[TestLogin] Полученный токен: %v\n", token)
	token, err = Login("admin", "")
	if (token != "") || (err == nil) {
		t.Error("[TestLogin] Не обрабатывается неправильный пароль")
	}
	token, err = Login("", "")
	if (token != "") || (err == nil) {
		t.Error("[TestLogin] Не обрабатывается неправильное имя пользователя")
	}
}
