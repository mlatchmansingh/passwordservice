package repository

import "PasswordService/domain/entities"

type InMemDB struct {
	Passwords map[int64]entities.Password
}
