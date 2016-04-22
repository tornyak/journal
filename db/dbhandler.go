///////////////////////////////////////////////////////////////////////////////
// package db provides interface to the DB.
// Communicaton is not direct but it is using GORP package providing ORM
// As a DB it is using sqlite3
///////////////////////////////////////////////////////////////////////////////

package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
	"fmt"
	"log"
)

const(
	DBDefaultLocation = "/tmp/quinn_db.bin"
	TableNameJournal = "journal"
)

// Column names
const(
	ColumnNameID = "id"
	ColumnNameName = "name"
	ColumnNameDuration = "duration"
	ColumnNameReason = "reason"
)

// Queries
var(
	QueryList = fmt.Sprintf("select * from %s order by %s", TableNameJournal, ColumnNameID)
	QueryTotal = fmt.Sprintf("select sum(%s) from %s", ColumnNameDuration, TableNameJournal)
	QueryHitList = fmt.Sprintf("select %s, sum(%s) from %s group by name order by sum(%s) desc",
		ColumnNameName, ColumnNameDuration, TableNameJournal, ColumnNameDuration)
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

// Hit is structure used for hitlist query
type Hit struct {
	Name string
	Duration int64 `db:"sum(duration)"`
}

// NewDBHandler will open database and create jorunal table if it does not exist
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
	checkErr(err, "Log: Insert failed")
}

// List all Interrupts from the DB
func (db *DBHandler) List() []Interrupt {
	// fetch all rows
	var interrupts []Interrupt
	_, err := db.dbMap.Select(&interrupts, QueryList)
	checkErr(err, "List: Select failed")
	return interrupts
}

// Total sums up and returns all interrupted time
func (db *DBHandler) Total() int64 {
	total, err := db.dbMap.SelectInt(QueryTotal)
	checkErr(err, "Total: Select failed")
	return total
}

// Total sums up and returns all interrupted time
func (db *DBHandler) Hitlist() []Hit {
	var hits []Hit
	_, err := db.dbMap.Select(&hits, QueryHitList)
	checkErr(err, "HitList: Select failed")
	return hits
}



func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	db, err := sql.Open("sqlite3", DBDefaultLocation)
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'journal' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Interrupt{}, TableNameJournal).SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
