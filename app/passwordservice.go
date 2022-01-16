package app

import (
	"PasswordService/domain/entities"
	"PasswordService/infrastructure/repository"
	"crypto/sha512"
	"encoding/base64"
	"log"
	"time"
)

type PasswordService interface {
	CreatePassword(password PasswordDTO) (PasswordEntityIdDTO, error)
	GetPassword(id entities.ID) (entities.Password, error)
	UpdatePassword(password entities.Password) (entities.ID, error)
}

type passwordService struct {
	repo *repository.PasswordRepository
}

func NewPasswordService(r *repository.PasswordRepository) PasswordService {
	return &passwordService{
		repo: r,
	}
}

func (s *passwordService) CreatePassword(plaintext PasswordDTO) (PasswordEntityIdDTO, error) {
	id, err := s.repo.Create(entities.NewPassword(plaintext.Password))
	if err != nil {
		return PasswordEntityIdDTO{}, err
	}

	time.AfterFunc(time.Second*5, func() { s.generatePasswordHash(id) })

	return PasswordEntityIdDTO{
		ID: id,
	}, nil
}

func (s *passwordService) GetPassword(id entities.ID) (entities.Password, error) {
	return s.repo.QueryById(id)
}

func (s *passwordService) UpdatePassword(password entities.Password) (entities.ID, error) {
	return s.repo.Update(password)
}

func (s *passwordService) generatePasswordHash(id entities.ID) {
	log.Printf("Generating hash for id %d\n", id)
	p, err := s.repo.QueryById(id)
	if err != nil {
		log.Fatalf("could not find a password entry for id %d\n", id)
	}

	h512 := sha512.New()
	h512.Write([]byte(p.Password))
	var hashed = h512.Sum(nil)
	var b64enc = base64.URLEncoding.EncodeToString(hashed)
	p.Password = b64enc
	p.Converted = true
	s.repo.Update(p)
}
