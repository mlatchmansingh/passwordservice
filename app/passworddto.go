package app

import "PasswordService/domain/entities"

type PasswordDTO struct {
	Password string `json:"password"`
}

type PasswordEntityIdDTO struct {
	ID entities.ID `json:"id"`
}

type PasswordStatsDTO struct {
	NumPosts int64 `json:"total"`
	Average  int64 `json:"average"`
}
