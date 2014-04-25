package zhDicts

import (
    "os"
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

func LoadDb() error {
    var err error
    dbname := os.Getenv("BMT_DB")
    dbuser := os.Getenv("BMT_USER")
    dbpw := os.Getenv("BMT_PW")

	db, err = sql.Open("postgres", "user=" + dbuser + " dbname=" + dbname + " password=" + dbpw + " host=localhost sslmode=disable")
    if err != nil {
        return err
    }

    return nil
}

func CloseDb() {
    db.Close()
}

func Db() *sql.DB {
    return db
}
