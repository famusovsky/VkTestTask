package models

// User - структура, представляющая пользователя.
type User struct {
	Id       int    `json:"id" db:"id"`             // Id - id пользователя.
	Nickname string `json:"nickname" db:"nickname"` // Nickname - никнейм (логин) пользователя.
	Password string `json:"password db:"password"`  // Password - пароль пользователя.
	IsAdmin  bool   `json:"is_admin" db:"is_admin"` // IsAdmin - флаг, указывающий на то, является ли пользователь администратором.
}
