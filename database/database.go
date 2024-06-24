package database

import (
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type DB struct {
	*pg.DB
}

func InitDB() (*DB, error) {
	options, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	pgDB := pg.Connect(options)

	_, err = pgDB.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	db := &DB{pgDB}

	err = createSchema(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *DB) error {
	models := []interface{}{
		(*User)(nil),
		(*UserProfile)(nil),
		(*Socials)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
