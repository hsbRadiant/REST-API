package configuration

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

// WAY TO CONFIGURE - (taught in Golang databse access tutorial) :
func NewDatabase() *sql.DB {
	cnfg := mysql.Config{
		User:      "root",          // os.Getenv("DBUSER")
		Passwd:    "NoorBedi@1997", //  os.Getenv("DBPASS")
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "test",
		ParseTime: true, // CHECK? (Helps to scan Date and Datetime to time.Time - READ more.)
	}

	db, err := sql.Open("mysql", cnfg.FormatDSN()) // 'FormatDSN()' formats the given 'Config' into a DSN string which can then be passed to the driver.
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	/*
		cntx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10) // READ? / CHECK?
		defer cancelFunc()
		db.ExecContext(cntx, "CREATE DATABASE IF NOT EXISTS "+cnfg.DBName)

		// _, err = db.Exec("CREATE SCHEMA IF NOT EXISTS ?", cnfg.DBName)

		if err != nil {
			log.Fatal(err.Error())
			return nil
		}
	*/

	db.SetConnMaxLifetime(time.Minute * 5) // CHECK?
	return db
}

var DB *sql.DB = NewDatabase() // a global variable to store the returned 'db' value of 'NewDatabase' function. Also it will help other files to store the connection

/*
func NewDatabase() *sql.DB {
	dbname := "test"
	user := "root"
	password := "NoorBedi@1997"
	address := "127.0.0.1:3306"

	// dataSource := user + ":" + password + "@tcp(" + address + ")/" + dbname
	// The 'Open' method does not creates a connection but only verify it's arg. (Hover over 'Open' method)
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+address+")/"+dbname)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Successfully connected to the database")
	}
	return db
}
*/
