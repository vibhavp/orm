package orm

import (
	"database/sql"
	"reflect"
)

type DB struct {
	DB *sql.DB
}

//Open opens a connection to the provided database, similar to database/sql's Open().
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return &DB{db}, err
}

//Migrate creates a table for the provided struct or struct pointer with the
//struct's name as the table name. It panics if v's is not a struct or a struct
//pointer.
func (db *DB) Migrate(v interface{}) error {
	return db.MigrateName(v, reflect.TypeOf(v).Name())
}

//MigrateName is like Migrate, but uses the provided name for the struct instead
//of the name of the struct's type.
func (db *DB) MigrateName(v interface{}, name string) error {
	var value = reflect.ValueOf(v)

	if reflect.TypeOf(v).Kind() != reflect.Struct {
		if reflect.TypeOf(v).Kind() == reflect.Ptr {
			if reflect.ValueOf(v).Elem().Kind() == reflect.Struct {
				value = reflect.ValueOf(v)
			}
		} else {
			panic("orm: call of Migrate on " + value.Kind().String())
		}
	}

	query := "CREATE TABLE " + name + "("
	if query == "" {
		panic("orm: call of Migrate on an unnamed type")
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)

		if field.Tag.Get("sql") == "-" {
			continue
		}

		switch value.Field(i).Interface().(type) {
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
	_, err := db.DB.Exec(query)

	return err

}
