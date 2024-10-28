package pg

import "github.com/jmoiron/sqlx"

func New(DSN string) (*DB, error) {
	db, err := sqlx.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	err = migrate(db)
	if err != nil {
		return nil, err
	}
	return &DB{
		con: db,
	}, nil
}
