package model

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/huandu/go-sqlbuilder"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Legitt@123"
	db_name  = "legittpro"
)

// Thread safe db pooled connection instance
var pool *sql.DB

func InitDb() *sql.DB {
	const dataSourceName = "postgres://postgres:Legitt@123@localhost:5432/legittpro?sslmode=disable"

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal("Unable initialise driver check data source string again", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	pool = db
	fmt.Print("Connected to database successfully\n")

	return pool
}

func InsertIntoDb(tableName string, columns []string, values [][]interface{}) (sql.Result, error) {
	sb := sqlbuilder.PostgreSQL.NewInsertBuilder().
		InsertInto(tableName).
		Cols(columns...)

	for _, val := range values {
		sb.Values(val...)
	}

	sql, args := sb.Build()

	res, err := pool.Exec(sql, args...)
	return res, err
}

func GetRowFromDb(tableName string, whereFields map[string]interface{}, fields ...string) *sql.Row {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select(fields...)
	sb.From(tableName)

	var whereClause []string
	for key, val := range whereFields {
		whereClause = append(whereClause, sb.Equal(key, val))
	}

	query, args := sb.Where(whereClause...).Build()

	row := pool.QueryRow(query, args...)

	return row
}

func GetRowsFromDb(tableName string, fields []string, whereClause string) (*sql.Rows, error) {
	queryString, args := sqlbuilder.
		Select(fields...).
		From(tableName).
		Where(whereClause).
		Build()

	rows, err := pool.Query(queryString, args...)

	return rows, err
}

func UpdateRowsInDb(tableName string, fieldsToUpdate []string, valuesToUpdate []interface{}, whereClause map[string]interface{}) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update(tableName)
	for index, element := range fieldsToUpdate {
		ub.SetMore(ub.Assign(element, valuesToUpdate[index]))
	}

	var whereClauseString []string
	for key, val := range whereClause {
		whereClauseString = append(whereClauseString, ub.Equal(key, val))
	}

	ub.Where(whereClauseString...)

	query, args := ub.Build()

	_, err := pool.Exec(query, args...)
	fmt.Print(err)

	return err
}
