package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DB struct {
	*sql.DB
}

// "user:password@/dbname"
func NewConnection(dsn string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

type Data struct {
	ID     int64  `json:"id"`
	Source string `json:"source"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Value  string `json:"value"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (data *Data) Scan(rows *sql.Rows) error {
	if err := rows.Scan(
		&data.ID, &data.Source, &data.Name, &data.Type, &data.Value, &data.CreatedAt, &data.UpdatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (db *DB) SaveData(item *Data) (int64, error) {
	// Prepare statement for inserting data
	stmt, err := db.Prepare("INSERT INTO data( source, name, type, value ) VALUES( ?, ?, ? )")
	if err != nil {
		panic(err.Error()) // TODO: proper error handling instead of panic in your app
	}
	defer stmt.Close() // Close the statement when we leave main() / the program terminates

	res, err := stmt.Exec(item.Source, item.Name, item.Type, item.Value)
	if err != nil {
		panic(err.Error()) // TODO: proper error handling instead of panic in your app
	}

	return res.LastInsertId()
}

func (db *DB) GetData() ([]*Data, error) {
	// Prepare statement for reading data
	stmt, err := db.Prepare("SELECT id, source, name, type, value, created_at, updated_at FROM data")
	if err != nil {
		panic(err.Error()) // TODO: proper error handling instead of panic in your app
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		panic(err.Error()) // TODO: proper error handling instead of panic in your app
	}
	defer rows.Close()

	var results []*Data
	for rows.Next() {
		result := &Data{}
		if err := result.Scan(rows); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
