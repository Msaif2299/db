package main

import (
	"fmt"
	"os"
	"time"

	"log/slog"

	MySQL "github.com/SpaceTent/db/mysql"
)

type UpdatePerson struct {
	Id      int       `db:"column=id primarykey=yes table=Users"`
	Name    string    `db:"column=name"`
	Dtadded time.Time `db:"column=dtadded"`
	Status  int       `db:"column=status"`
	Ignored int       `db:"column=ignored omit=yes"`
}

type InsertPerson struct {
	Id      int       `db:"column=id primarykey=yes table=Users"`
	Name    string    `db:"column=name"`
	Dtadded time.Time `db:"column=dtadded"`
	Status  int       `db:"column=status"`
	Ignored int       `db:"column=ignored omit=yes"`
}

// InsertOneExample shows how to insert a single row into the table
func InsertOneExample(l *slog.Logger) {
	// First create the structure
	entry := InsertPerson{
		Id:      12,
		Name:    "Test",
		Dtadded: time.Now(),
		Status:  1,
	}
	// Now create the query
	sqlQuery, err := MySQL.DB.Insert(entry)
	if err != nil {
		l.Error(err.Error())
		return
	}
	// Then execute the query
	lastInsertedID, rowsAffected, err := MySQL.DB.Execute(sqlQuery)
	if err != nil {
		l.Error(err.Error())
	}
	l.Info(fmt.Sprintf("Item with ID %d was inserted. %d rows were affected", lastInsertedID, rowsAffected))
}

// UpdateExample shows how to update an entry in the table
func UpdateExample(l *slog.Logger) {
	// First create the structure
	p := UpdatePerson{
		Id:      12,
		Name:    "Test",
		Dtadded: time.Now(),
		Status:  1,
	}
	// Now create the query
	sqlQuery, err := MySQL.DB.Update(p)
	if err != nil {
		l.Error(err.Error())
	}
	// Then execute it
	lastInsertedID, rowsAffected, err := MySQL.DB.Execute(sqlQuery)
	if err != nil {
		l.Error(err.Error())
		return
	}
	l.Info(fmt.Sprintf("Item with ID %d was updated. %d rows were affected", lastInsertedID, rowsAffected))
}

// HandleMultipleConnectionsExample shows how to handle multiple connections to different databases
func HandleMultipleConnectionsExample(l *slog.Logger) {
	type Connections struct {
		DatabaseA *MySQL.Database
		DatabaseB *MySQL.Database
	}
	conns := Connections{}
	MySQL.New("testA", l)
	conns.DatabaseA = MySQL.DB
	MySQL.New("testB", l)
	conns.DatabaseB = MySQL.DB

	type ExampleSelect struct {
		Id      int       `db:"column=id primarykey=yes table=Users"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded"`
	}

	result, err := MySQL.QuerySingleStructV2[ExampleSelect](conns.DatabaseA, "SELECT * FROM Users WHERE id=?", 5)
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Result: %+v", result))

	results, err := MySQL.QueryStructV2[ExampleSelect](conns.DatabaseB, "SELECT * FROM Users")
	if err != nil {
		return
	}
	for _, r := range results {
		l.Info(fmt.Sprintf("Result: %+v\n", r))
	}
}

func main() {

	DSN := ""
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	l := slog.New(textHandler)

	MySQL.New(DSN, l)
	InsertOneExample(l)
	UpdateExample(l)
}
