package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

const (
	DEFAULT_IMPLANT_COMM_PROTOCOL = HTTPS
)

// TODO: where to put this?
func setupDB() *gorm.DB {
	// setup DB
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable" // connect w/: `psql -h localhost -p 5432 -U postgres -d postgres`
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// TODO: how to read these logs?
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err = db.AutoMigrate(&Implant{}, &Task{}, &Result{}); err != nil {
		panic(err) // TODO: fix
	}
	return db
}

func main() {
	db := setupDB()

	httpsListener, _ := NewListener(db, HTTPS)
	httpsListener.listen()

	dnsListener, _ := NewListener(db, DNS)
	dnsListener.listen()
}
