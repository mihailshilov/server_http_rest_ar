package sqlstore

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/store"
)

//Stores
type Store struct {
	dbPostgres *pgxpool.Pool
	// dbMssql        *sql.DB
	userRepository *UserRepository
	dataRepository *DataRepository
}

//New_db
//func New(db *pgxpool.Pool, dbmssql *sql.DB) *Store {
func New(db *pgxpool.Pool) *Store {
	return &Store{
		dbPostgres: db,
		// dbMssql:    dbmssql,
	}
}

//User
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

//Data
func (s *Store) Data() store.DataRepository {
	if s.dataRepository != nil {
		return s.dataRepository
	}

	s.dataRepository = &DataRepository{
		store: s,
	}

	return s.dataRepository
}
