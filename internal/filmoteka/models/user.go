package models

import "errors"

// User - структура, представляющая пользователя.
type User struct {
	Id       int    `json:"-" db:"id"`              // Id - id пользователя.
	Nickname string `json:"nickname" db:"nickname"` // Nickname - никнейм (логин) пользователя.
	Password string `json:"password" db:"password"` // Password - пароль пользователя.
	IsAdmin  bool   `json:"is_admin" db:"is_admin"` // IsAdmin - флаг, указывающий на то, является ли пользователь администратором.
}

// Check - проверка корректности данных пользователя.
//
// Возвращает: ошибку.
func (u *User) Check() error {
	errs := make([]error, 0, 3)
	if u.Nickname == "" {
		errs = append(errs, errors.New("name must not be null"))
	}
	if u.Password == "" {
		errs = append(errs, errors.New("password must not be null"))
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}
