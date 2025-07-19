package model

import "time"

type CodewarsUser struct {
	Username  string    `json:"username"`
	Honor     int       `json:"honor"`
	CreatedAt time.Time // Для хранения в БД
}

type User struct {
	CodewarsUser
	CreatedAt time.Time
}
