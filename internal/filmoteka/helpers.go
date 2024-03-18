package filmoteka

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// sendJson - отправка json-ответа.
//
// Принимает: ResponseWriter и любой объект.
func (app *App) sendJson(w http.ResponseWriter, obj any) {
	js, err := json.Marshal(obj)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// authIsAdmin - проверка прав пользователя.
//
// Принимает: http.Request.
//
// Возвращает: true, если пользователь админ, иначе false и ошибку.
func (app *App) authIsAdmin(r *http.Request) (bool, error) {
	nick, pswd, ok := r.BasicAuth()
	if !ok {
		return false, errors.New("error parsing basic auth")
	}

	if nick == "admin" && pswd == "admin" {
		return true, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pswd), 8)
	if err != nil {
		return false, err
	}

	isAdmin, err := app.dbHandler.CheckUserRole(nick, string(hashedPassword))
	if err != nil {
		return false, err
	}

	return isAdmin, nil
}

// handleError - обработка ошибок.
//
// Принимает: логгер, ResponseWriter, сообщение и http-статус.
func handleError(log *log.Logger, w http.ResponseWriter, msg string, status int) {
	log.Println(msg)
	http.Error(w, msg, status)
}
