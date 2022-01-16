package app

import (
	"PasswordService/domain/entities"
	"PasswordService/infrastructure/repository"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"log"
	"sync"
	"time"
)

type PasswordService interface {
	CreatePassword(password PasswordDTO) (PasswordEntityIdDTO, error)
	GetPassword(id entities.ID) (entities.Password, error)
	UpdatePassword(password entities.Password) (entities.ID, error)
	CloseAndWait()
}

type removeLock struct {
	removed bool
	wg      sync.WaitGroup
	lock    sync.Mutex
}

type passwordService struct {
	repo   *repository.PasswordRepository
	iolock removeLock
}

func NewPasswordService(r *repository.PasswordRepository) PasswordService {
	ps := &passwordService{
		repo: r,
	}

	ps.iolock.wg.Add(1)
	ps.iolock.removed = false
	return ps
}

func (s *passwordService) acquireRemoveLock() error {
	s.iolock.lock.Lock()
	defer s.iolock.lock.Unlock()
	if s.iolock.removed {
		return errors.New("not accepting new requests")
	}
	s.iolock.wg.Add(1)
	return nil
}

func (s *passwordService) releaseRemoveLock() {
	s.iolock.wg.Done()
}

func (s *passwordService) releaseAndWaitRemoveLock() {
	s.iolock.lock.Lock()
	s.iolock.removed = true
	s.iolock.lock.Unlock()

	//
	// decrement by 1 for the lock we took in Init
	s.iolock.wg.Done()
	s.iolock.wg.Wait()
}

func (s *passwordService) CreatePassword(plaintext PasswordDTO) (PasswordEntityIdDTO, error) {

	err := s.acquireRemoveLock()
	if err != nil {
		return PasswordEntityIdDTO{}, err
	}

	id, err := s.repo.Create(entities.NewPassword(plaintext.Password))
	if err != nil {
		s.releaseRemoveLock()
		return PasswordEntityIdDTO{}, err
	}

	time.AfterFunc(time.Second*5, func() { s.generatePasswordHash(id) })

	return PasswordEntityIdDTO{
		ID: id,
	}, nil
}

func (s *passwordService) GetPassword(id entities.ID) (entities.Password, error) {
	err := s.acquireRemoveLock()
	if err != nil {
		return entities.Password{}, err
	}

	p, err := s.repo.QueryById(id)
	s.releaseRemoveLock()
	return p, err
}

func (s *passwordService) UpdatePassword(password entities.Password) (entities.ID, error) {
	err := s.acquireRemoveLock()
	if err != nil {
		return entities.InvalidID, err
	}
	id, err := s.repo.Update(password)
	s.releaseRemoveLock()
	return id, err
}

func (s *passwordService) CloseAndWait() {

	s.releaseAndWaitRemoveLock()

}

func (s *passwordService) generatePasswordHash(id entities.ID) {
	log.Printf("Generating hash for id %d\n", id)
	p, err := s.repo.QueryById(id)
	if err != nil {
		s.releaseRemoveLock()
		log.Fatalf("could not find a password entry for id %d\n", id)
	}

	h512 := sha512.New()
	h512.Write([]byte(p.Password))
	var hashed = h512.Sum(nil)
	var b64enc = base64.URLEncoding.EncodeToString(hashed)
	p.Password = b64enc
	p.Converted = true
	s.repo.Update(p)
	s.releaseRemoveLock()
}
