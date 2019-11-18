package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dasider41/cfootsell/util"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

// InitDB :
func InitDB() *sql.DB {
	env, err := godotenv.Read()
	util.ErrCheck(err)

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		env["DB_USERNAME"],
		env["DB_PASSWORD"],
		env["DB_HOST"],
		env["DB_PORT"],
		env["DB_DATABASE"])
	db, err := sql.Open("mysql", dbConn)
	util.ErrCheck(err)
	return db
}

// DateFormat :
func DateFormat(tDate string) string {
	layoutIN := "06-01-02"
	layoutOUT := "2006-01-02"
	t, err := time.Parse(layoutIN, tDate)

	if err != nil {
		return time.Now().Format(layoutOUT)
	}

	return t.Format(layoutOUT)
}
