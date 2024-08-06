package api

import (
	"sync"
)

type DB struct {
	data sync.Map
}

func NewDatabase() *DB {
	return &DB{}
}

func (db *DB) Set(username string, userData User) {
	db.data.Store(username, userData)
}

func (db *DB) Get(username string) (User, bool) {
	res, ok := db.data.Load(username)
	if !ok {
		return User{}, false
	}

	user, ok := res.(User)
	if !ok {
		return User{}, false
	}

	return user, true
}

func (db *DB) Delete(username string) {
	db.data.Delete(username)
}
