package db

import (
	"fmt"
	"strings"

	"sync"
	"test/internal/error"
	"test/internal/utils"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var mu sync.Mutex

// Assuming currentDB.PrimaryKey is defined somewhere in your code
var currentDB = struct {
	PrimaryKey string
}{
	PrimaryKey: "INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE",
}

// Connect to SQLite database
func ConnectDB() (*sqlx.DB, bool) {
	// allHosts []utils.Host
	fmt.Println("connectdb")
	var ok bool
	db, err := sqlx.Open("sqlite", "file:/home/user/Desktop/workspace/grafana/grafana-NMS/device/device.db?mode=rw")
	error.IfError(err)

	err = db.Ping()

	if error.IfError(err) && db.DriverName() != "sqlite" {

		Create()
	}

	// error.IfError(err)
	ok = true

	return db, ok
}

// Execute a given SQL statement
func dbExec(sqlStatement string) {
	fmt.Println("dbexec")

	db, ok := ConnectDB()
	defer db.Close()

	if ok {
		mu.Lock()
		_, err := db.Exec(sqlStatement)
		mu.Unlock()
		error.IfError(err)
	}
}

// Create tables "now" and "history"
func Create() {
	sqlStatement := `CREATE TABLE IF NOT EXISTS "now" (
		"ID"	` + currentDB.PrimaryKey + `,
		"NAME"	TEXT,
		"DNS"   TEXT,
		"IFACE"	TEXT,
		"IP"	TEXT,
		"MAC"	TEXT,
		"HW"	TEXT,
		"DATE"	TEXT,
		"KNOWN"	INTEGER DEFAULT 0,
		"NOW"	INTEGER DEFAULT 0
	);`
	dbExec(sqlStatement)

	sqlStatement = `CREATE TABLE IF NOT EXISTS "history" (
		"ID"	` + currentDB.PrimaryKey + `,
		"NAME"	TEXT,
		"DNS"	TEXT,
		"IFACE"	TEXT,
		"IP"	TEXT,
		"MAC"	TEXT,
		"HW"	TEXT,
		"DATE"	TEXT,
		"KNOWN"	INTEGER DEFAULT 0,
		"NOW"	INTEGER DEFAULT 0
	);`
	dbExec(sqlStatement)
}

// Insert data into the "now" table
func Insert(table string, oneHost utils.Host) {
	fmt.Println("inserting in to database________________________")
	// oneHost.Name = quoteStr(oneHost.Name)
	// oneHost.Hw = quoteStr(oneHost.Hw)
	sqlStatement := `INSERT INTO %s ("NAME", "DNS", "IFACE", "IP", "MAC", "HW", "DATE", "KNOWN", "NOW") VALUES ('%s','%s','%s','%s','%s','%s','%s','%d','%d');`
	// fmt.Println("SQL Statement:", sqlStatement)
	sqlStatement = fmt.Sprintf(sqlStatement, table, oneHost.Name, oneHost.DNS, oneHost.Iface, oneHost.IP, oneHost.Mac, oneHost.Hw, oneHost.Date, oneHost.Known, oneHost.Now)
	// fmt.Println("SQL Statement:", sqlStatement)
	fmt.Printf("%s\n", oneHost.Name)
	fmt.Printf("%s\n", oneHost.Iface)
	fmt.Printf("%s\n", oneHost.IP)

	dbExec(sqlStatement)
	fmt.Println("end inserting in to database____________")
	// mu.Unlock()

}
func Select(table string) (dbHosts []utils.Host) {

	sqlStatement := `SELECT * FROM ` + table + ` ORDER BY "DATE" DESC`

	db, ok := ConnectDB()
	defer db.Close()

	if ok {
		mu.Lock()
		err := db.Select(&dbHosts, sqlStatement)
		mu.Unlock()
		error.IfError(err)
	}

	return dbHosts
}

func Update(table string, oneHost utils.Host) {
	// fmt.Println("updating table___________________________________________")
	// oneHost.Name = quoteStr(oneHost.Name)
	// oneHost.Hw = quoteStr(oneHost.Hw)
	sqlStatement := `UPDATE %s set
		"NAME" = '%s', "DNS" = '%s', "IFACE" = '%s', "IP" = '%s', "MAC" = '%s', "HW" = '%s', "DATE" = '%s', "KNOWN" = '%d', "NOW" = '%d'
		WHERE "ID" = '%d';`

	sqlStatement = fmt.Sprintf(sqlStatement, table, oneHost.Name, oneHost.DNS, oneHost.Iface, oneHost.IP, oneHost.Mac, oneHost.Hw, oneHost.Date, oneHost.Known, oneHost.Now, oneHost.ID)

	dbExec(sqlStatement)
	// fmt.Println("end table___________________________________________")
}
func Delete(table string, id int) {
	sqlStatement := `DELETE FROM %s WHERE "ID"='%d';`
	sqlStatement = fmt.Sprintf(sqlStatement, table, id)
	dbExec(sqlStatement)
}

func quoteStr(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}
