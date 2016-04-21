package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
	"log"
)

// DBHandler provides DB interface for journal CLI commands
type DBHandler struct {
	dbMap *gorp.DbMap
}

// Interrupt will be mapped into a DB table row
type Interrupt struct {
	Id       int64
	Name     string
	Duration int64
	Reason   string
}

type Hit struct {
	Name string
	Duration int64 `db:"sum(duration)"`
}

func NewDBHandler() *DBHandler {
	return &DBHandler{
		dbMap: initDb(),
	}
}

// Insert a new Interrupt into the DB
func (db *DBHandler) Log(name string, duration int64, reason string) {
	interrupt := &Interrupt{
		Name:     name,
		Duration: duration,
		Reason:   reason,
	}
	err := db.dbMap.Insert(interrupt)
	checkErr(err, "Insert failed")
}

// List all Interrupts from the DB
func (db *DBHandler) List() []Interrupt {
	// fetch all rows
	var interrupts []Interrupt
	_, err := db.dbMap.Select(&interrupts, "select * from journal order by id")
	checkErr(err, "Select failed")
	return interrupts
}

// Total sums up and returns all interrupted time
func (db *DBHandler) Total() int64 {
	total, err := db.dbMap.SelectInt("select sum(duration) from journal")
	checkErr(err, "Select failed")
	return total
}

// Total sums up and returns all interrupted time
func (db *DBHandler) Hitlist() []Hit {
	var hits []Hit
	_, err := db.dbMap.Select(&hits, "select name, sum(duration) from journal group by name order by sum(duration) desc")
	checkErr(err, "Select failed")
	return hits
}



func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	db, err := sql.Open("sqlite3", "/tmp/quinn_db.bin")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'journal' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Interrupt{}, "journal").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
