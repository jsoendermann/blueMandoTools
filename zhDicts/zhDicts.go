package zhDicts

import (
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

func LoadDb() error {
    var err error
    db, err = sql.Open("postgres", "user=json dbname=json host=localhost sslmode=disable")
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
