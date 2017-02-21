//model package uses Bolt DB connection for interacting with store
package model

import (
	"fmt"

	"github.com/alexivanenko/egroupware_notifier_bot/config"
	"github.com/boltdb/bolt"
)

var db *bolt.DB

func init() {
	dbPath := fmt.Sprintf("%s/%s", config.GetRootDir(), config.String("db", "store"))

	config.Log(fmt.Sprintf("DB store file: %s", dbPath))

	db, _ = bolt.Open(dbPath, 0644, nil)
}

//GetDB returns a pointer to the database
func GetDB() *bolt.DB {
	return db
}
