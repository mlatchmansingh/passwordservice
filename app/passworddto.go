package app

import "PasswordService/domain/entities"

type PasswordDTO struct {
	Password string `json:"password"`
}

type PasswordEntityIdDTO struct {
	ID entities.ID `json:"id"`
}
