package repository

import (
	"PasswordService/domain/entities"
	"errors"
	"sync"
)

//
// The repository layer will implement a database in memory using a dictionary.
// The interface is abstracted so that the outer layers don't have knowledge of how
// data is saved. It could be file, database, or memory
type PasswordRepository struct {
	Db   *InMemDB
	lock sync.RWMutex
}

func (p *PasswordRepository) Init() {
	p.Db.Passwords = make(map[entities.ID]entities.Password)
}

func (p *PasswordRepository) QueryById(id entities.ID) (entities.Password, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if value, ok := p.Db.Passwords[id]; ok {
		return value, nil
	}
	return entities.Password{}, errors.New("record not found")
}

func (p *PasswordRepository) Create(password entities.Password) (entities.ID, error) {
	//
	// We simulate using the password's Id field as the primary key
	// and we depend on the user to ensure it is unique
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, ok := p.Db.Passwords[password.Id]; ok {
		//
		// the primary key is not unique - fail
		return entities.InvalidID, errors.New("primary key is not unique")
	}

	p.Db.Passwords[password.Id] = password
	return password.Id, nil
}

func (p *PasswordRepository) Update(password entities.Password) (entities.ID, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, ok := p.Db.Passwords[password.Id]; !ok {
		return entities.InvalidID, errors.New("primary key does not exist")
	}
	p.Db.Passwords[password.Id] = password
	return password.Id, nil
}
