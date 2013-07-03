package zhDicts

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "os"
)

var db *sql.DB

func LoadDb() error {
    dbPath := os.Getenv("ZHDICT_DB")

    var err error
    db, err = sql.Open("sqlite3", dbPath)
    if err != nil {
        return err
    }

    return nil
}

func CloseDb() {
    db.Close()
}

func getDb() *sql.DB {
    return db
}
