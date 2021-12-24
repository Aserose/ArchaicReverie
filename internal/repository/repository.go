package repository

import "github.com/Aserose/ArchaicReverie/internal/repository/postgres/data"

type DB struct {
	Postgres *data.PostgresData
}

func NewDB(postgres *data.PostgresData) *DB {
	return &DB{
		Postgres: postgres,
	}
}
