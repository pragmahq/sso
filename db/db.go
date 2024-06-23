package db

import (
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func InitDB() (*pg.DB, error) {
	options, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	db := pg.Connect(options)

	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	err = createSchema(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *pg.DB) error {
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
