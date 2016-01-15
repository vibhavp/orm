package orm

import (
	"database/sql"
	"reflect"
)

//DB ...
type DB struct {
	db *sql.DB
}

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return &DB{db}, err
}

func (db *DB) AutoMigrate(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		panic("orm: argument not a struct")
	}

	query := "CREATE TABLE " + reflect.TypeOf(v).Name() + "("
	if query == "" {
		panic("orm: argument is an unnamed type")
	}

	for i := 0; i < reflect.TypeOf(v).NumField(); i++ {
		field := reflect.TypeOf(v).Field(i)

		if field.Tag.Get("sql") == "-" {
			continue
		}

		switch reflect.ValueOf(v).Field(i).Interface().(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			query += field.Name + " bigint"
		case string:
			query += field.Name + " varchar"
		case bool:
			query += field.Name + " bool"
		}
		if i == reflect.TypeOf(v).NumField()-1 {
			continue
		}
		query += ","

	}
	query += ")"
	_, err := db.db.Exec(query)

	return err
}
